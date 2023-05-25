package services

import (
	"bytes"
	"fmt"
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

// DefaultKnowledgeBaseName its default name of Knowledge Base
const DefaultKnowledgeBaseName = "default"

// DefaultKnowledgeBaseVersion its default version of Knowledge Base
const DefaultKnowledgeBaseVersion = "latest"

// LoadLocalGRL loads a GRL (Grule Rule Language) file from a local path and builds a rule from it
// using the `builder.NewRuleBuilder` function. It takes in the path of the GRL file, the name of the
// knowledge base, and the version of the knowledge base as parameters. It returns an error if there is
// any issue building the rule from the resource.
func (s Eval) LoadLocalGRL(grlPath string, knowledgeBaseName string, version string) error {
	ruleBuilder := builder.NewRuleBuilder(s.knowledgeLibrary)
	fileRes := pkg.NewFileResource(grlPath)
	return ruleBuilder.BuildRuleFromResource(knowledgeBaseName, version, fileRes)
}

// The type `knowledgeBaseInfo` contains information about a knowledge base (rule sheet), including its name and version.
//
// Property
//   - KnowledgeBaseName: The property is a string that represents the name of a knowledge base. A knowledge base is a repository of information used to support decision-making, problem-solving, and other activities. In this case, each knowledge base is a rule sheet and contains the rules within it.
//   - Version: It's a string that represents the version number of the rule sheet, that is, the version of the knowledge base you want to use.
type knowledgeBaseInfo struct {
	KnowledgeBaseName string
	Version           string
}

// LoadRemoteGRL function is responsible for loading GRL (Grule Rule Language) rules from a remote location, such as a GitLab repository,
// and constructing a rule from them using the builder.NewRuleBuilder function. It takes the knowledge base name (rulesheet) and the
// knowledge base version as parameters.
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

// evalMutex is a variable of type sync.Mutex. This variable is used to synchronize access to the Eval method of the Eval struct, which is concurrently
// called by multiple goroutines. By using a mutex, it ensures that only one goroutine can execute the Eval method at a time, preventing race conditions and ensuring proper execution of the method.
var evalMutex sync.Mutex

// IEval ...
type IEval interface {
	GetKnowledgeLibrary() *ast.KnowledgeLibrary
	GetDefaultKnowledgeBase() *ast.KnowledgeBase
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

// GetDefaultKnowledgeBase ...
func (s Eval) GetDefaultKnowledgeBase() *ast.KnowledgeBase {
	return s.GetKnowledgeLibrary().GetKnowledgeBase(DefaultKnowledgeBaseName, DefaultKnowledgeBaseVersion)
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
