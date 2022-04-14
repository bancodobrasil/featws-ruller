package services

import (
	"bytes"
	"strings"
	"sync"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"

	"github.com/bancodobrasil/featws-ruller/config"
	"github.com/bancodobrasil/featws-ruller/processor"
	"github.com/bancodobrasil/featws-ruller/types"
)

//LoadLocalGRL ...
func (s Eval) LoadLocalGRL(grlPath string, knowledgeBaseName string, version string) error {
	ruleBuilder := builder.NewRuleBuilder(s.knowledgeLibrary)
	fileRes := pkg.NewFileResource(grlPath)
	return ruleBuilder.BuildRuleFromResource(knowledgeBaseName, version, fileRes)
}

type knowledgeBaseInfo struct {
	KnowledgeBaseName string
	Version           string
}

//LoadRemoteGRL ...
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
	LoadLocalGRL(grlPath string, knowledgeBaseName string, version string) error
	LoadRemoteGRL(knowledgeBaseName string, version string) error
	Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (*types.Result, error)
}

// EvalService ...
var EvalService IEval = NewEval()

// Eval ... struct
type Eval struct {
	knowledgeLibrary *ast.KnowledgeLibrary
}

// NewEval ...
func NewEval() Eval {
	return Eval{
		knowledgeLibrary: ast.NewKnowledgeLibrary(),
	}
}

// GetKnowledgeLibrary ...
func (s Eval) GetKnowledgeLibrary() *ast.KnowledgeLibrary {
	return s.knowledgeLibrary
}

//Eval ...
func (s Eval) Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (*types.Result, error) {
	// FIXME Remove synchronization on eval
	evalMutex.Lock()
	dataCtx := ast.NewDataContext()

	processor := processor.NewProcessor()

	result := types.NewResult()

	err := dataCtx.Add("processor", processor)
	if err != nil {
		log.Error("error on add processor to data context: \n the result was: %w \n the error was: %w", result, err)
		return result, err
	}

	err = dataCtx.Add("ctx", ctx)
	if err != nil {
		log.Error("error on add context to data context: \n the result was: %w \n the error was: %w", result, err)
		return result, err
	}

	err = dataCtx.Add("result", result)
	if err != nil {
		log.Error("error on add result to data context: \n the result was: %w \n the error was: %w", result, err)
		return result, err
	}

	eng := engine.NewGruleEngine()
	err = eng.Execute(dataCtx, knowledgeBase)
	if err != nil {
		log.Error("error on execute the grule engine: %w", err)
		return nil, err
	}

	if ctx.Has("errors") && len(ctx.GetMap("errors").GetEntries()) > 0 {
		result.Put("errors", ctx.GetMap("errors").GetEntries())
	}

	log.Debug("Context:\n\t", ctx.GetEntries(), "\n\n")
	log.Debug("Features:\n\t", result.GetFeatures(), "\n\n")

	evalMutex.Unlock()

	return result, nil
}
