package services

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"

	"github.com/bancodobrasil/featws-ruller/config"
	"github.com/bancodobrasil/featws-ruller/processor"
	"github.com/bancodobrasil/featws-ruller/types"
)

// DefaultKnowledgeBaseName its default name of Knowledge Base
const DefaultKnowledgeBaseName = "default"

// DefaultKnowledgeBaseVersion its default version of Knowledge Base
const DefaultKnowledgeBaseVersion = "latest"

// LoadLocalGRL ...
func (s Eval) LoadLocalGRL(grlPath string, knowledgeBaseName string, version string) error {
	ruleBuilder := builder.NewRuleBuilder(s.knowledgeLibrary)
	fileRes := pkg.NewFileResource(grlPath)
	return ruleBuilder.BuildRuleFromResource(knowledgeBaseName, version, fileRes)
}

type knowledgeBaseInfo struct {
	KnowledgeBaseName string
	Version           string
}

type knowledgeBaseCache struct {
	KnowledgeBase  *ast.KnowledgeBase
	ExpirationDate time.Time
}

// LoadRemoteGRL ...
func (s Eval) LoadRemoteGRL(knowledgeBaseName string, version string) error {
	cfg := config.GetConfig()
	ruleBuilder := builder.NewRuleBuilder(s.knowledgeLibrary)

	url := cfg.ResourceLoaderURL
	url = strings.Replace(url, "{knowledgeBase}", "{{.KnowledgeBaseName}}", -1)
	url = strings.Replace(url, "{version}", "{{.Version}}", -1)

	info := knowledgeBaseInfo{
		KnowledgeBaseName: knowledgeBaseName,
		Version:           version,
	}

	urlTemplate := template.New("UrlTemplate")

	// "Parse" parses a string into a template
	urlTemplate, _ = urlTemplate.Parse(url)

	var doc bytes.Buffer
	// standard output to print merged data
	err := urlTemplate.Execute(&doc, info)
	if err != nil {
		log.Error("error on load Remote GRL: %w", err)
		return err
	}

	url = doc.String()

	log.Debug("LoadRemoteGRL: ", url)

	fileRes := pkg.NewURLResourceWithHeaders(url, cfg.ResourceLoaderHeaders)
	return ruleBuilder.BuildRuleFromResource(knowledgeBaseName, version, fileRes)
}

var evalMutex sync.Mutex

// IEval ...
type IEval interface {
	GetKnowledgeLibrary() *ast.KnowledgeLibrary
	GetDefaultKnowledgeBase() *ast.KnowledgeBase
	GetKnowledgeBase(knowledgeBaseName string, version string) (*ast.KnowledgeBase, error)
	LoadLocalGRL(grlPath string, knowledgeBaseName string, version string) error
	LoadRemoteGRL(knowledgeBaseName string, version string) error
	Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (*types.Result, error)
}

var loadMutex sync.Mutex

// EvalService ...
var EvalService IEval = NewEval()

// Eval ... struct
type Eval struct {
	knowledgeLibrary   *ast.KnowledgeLibrary
	knowledgeBaseCache map[knowledgeBaseInfo]*knowledgeBaseCache
}

// NewEval ...
func NewEval() Eval {
	return Eval{
		knowledgeLibrary:   ast.NewKnowledgeLibrary(),
		knowledgeBaseCache: map[knowledgeBaseInfo]*knowledgeBaseCache{},
	}
}

// GetKnowledgeLibrary ...
func (s Eval) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return s.knowledgeLibrary
}

// GetDefaultKnowledgeBase ...
func (s Eval) GetDefaultKnowledgeBase() *ast.KnowledgeBase {

	return s.GetKnowledgeLibrary().GetKnowledgeBase(DefaultKnowledgeBaseName, DefaultKnowledgeBaseVersion)
}

// GetKnowledgeBase ...
func (s Eval) GetKnowledgeBase(knowledgeBaseName string, version string) (*ast.KnowledgeBase, error) {
	info := knowledgeBaseInfo{KnowledgeBaseName: knowledgeBaseName, Version: version}
	existing := s.knowledgeBaseCache[info]

	if existing == nil {

		existing = &knowledgeBaseCache{
			KnowledgeBase:  s.GetKnowledgeLibrary().GetKnowledgeBase(knowledgeBaseName, version),
			ExpirationDate: time.Now().Add(time.Minute * 5),
		}
		s.knowledgeBaseCache[info] = existing

	}

	if existing.ExpirationDate.After(time.Now()) || !(len(existing.KnowledgeBase.RuleEntries) > 0) {
		//invalidateCache
		loadMutex.Lock()

		err := s.LoadRemoteGRL(knowledgeBaseName, version)
		if err != nil {
			log.Errorf("Erro on load: %v", err)
			loadMutex.Unlock()
			return nil, errors.New("error on load knowledgebase and/or version")
		}

		if !(len(existing.KnowledgeBase.RuleEntries) > 0) {

			loadMutex.Unlock()
			return nil, errors.New("knowledgebase or version not found")
		}

		loadMutex.Unlock()

		existing.KnowledgeBase = s.GetKnowledgeLibrary().GetKnowledgeBase(knowledgeBaseName, version)
		existing.ExpirationDate = time.Now().Add(time.Minute * 5)

	}
	return existing.KnowledgeBase, nil

}

// Eval ...
func (s Eval) Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (result *types.Result, err error) {
	// FIXME Remove synchronization on eval
	evalMutex.Lock()
	defer evalMutex.Unlock()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
			log.Error(err)
		}
	}()
	dataCtx := ast.NewDataContext()

	processor := processor.NewProcessor()

	result = types.NewResult()

	err = dataCtx.Add("processor", processor)
	if err != nil {
		log.Error("error on add processor to data context: \n the result was: %w \n the error was: %w", result, err)
		return
	}

	err = dataCtx.Add("ctx", ctx)
	if err != nil {
		log.Error("error on add context to data context: \n the result was: %w \n the error was: %w", result, err)
		return
	}

	err = dataCtx.Add("result", result)
	if err != nil {
		log.Error("error on add result to data context: \n the result was: %w \n the error was: %w", result, err)
		return
	}

	eng := engine.NewGruleEngine()
	err = eng.Execute(dataCtx, knowledgeBase)
	if err != nil {
		log.Error("error on execute the grule engine: %w", err)
		return
	}

	if ctx.Has("errors") && len(ctx.GetMap("errors").GetEntries()) > 0 {
		result.Put("errors", ctx.GetMap("errors").GetEntries())
	}

	if ctx.Has("requiredParamErrors") && len(ctx.GetMap("requiredParamErrors").GetEntries()) > 0 {
		result.Put("requiredParamErrors", ctx.GetMap("requiredParamErrors").GetEntries())
	}

	log.Debug("Context:\n\t", ctx.GetEntries(), "\n\n")
	log.Debug("Features:\n\t", result.GetFeatures(), "\n\n")

	return
}
