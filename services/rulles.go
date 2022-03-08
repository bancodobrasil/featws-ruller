package services

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"sync"
	"text/template"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"

	"github.com/bancodobrasil/featws-ruller/config"
	"github.com/bancodobrasil/featws-ruller/processor"
	"github.com/bancodobrasil/featws-ruller/types"
)

//KnowledgeLibrary ...
var KnowledgeLibrary *ast.KnowledgeLibrary = ast.NewKnowledgeLibrary()

//LoadLocalGRL ...
func LoadLocalGRL(grlPath string, knowledgeBaseName string, version string) error {
	ruleBuilder := builder.NewRuleBuilder(KnowledgeLibrary)
	fileRes := pkg.NewFileResource(grlPath)
	return ruleBuilder.BuildRuleFromResource(knowledgeBaseName, version, fileRes)
}

type knowledgeBaseInfo struct {
	KnowledgeBaseName string
	Version           string
}

//LoadRemoteGRL ...
func LoadRemoteGRL(knowledgeBaseName string, version string) error {
	cfg := config.GetConfig()
	ruleBuilder := builder.NewRuleBuilder(KnowledgeLibrary)
	headers := make(http.Header)
	for header, value := range cfg.ResourceLoaderHeaders {
		headers.Set(header, value)
	}

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
		return err
	}

	url = doc.String()

	fileRes := pkg.NewURLResourceWithHeaders(url, headers)
	return ruleBuilder.BuildRuleFromResource(knowledgeBaseName, version, fileRes)
}

var evalMutex sync.Mutex

//Eval ...
func Eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (*types.Result, error) {
	// FIXME Remove synchronization on eval
	evalMutex.Lock()
	dataCtx := ast.NewDataContext()

	processor := processor.NewProcessor()

	result := types.NewResult()

	err := dataCtx.Add("processor", processor)
	if err != nil {
		return result, err
	}

	err = dataCtx.Add("ctx", ctx)
	if err != nil {
		return result, err
	}

	err = dataCtx.Add("result", result)
	if err != nil {
		return result, err
	}

	eng := engine.NewGruleEngine()
	err = eng.Execute(dataCtx, knowledgeBase)
	if err != nil {
		panic(err)
	}

	log.Print("Context:\n\t", ctx.GetEntries(), "\n\n")
	log.Print("Features:\n\t", result.GetFeatures(), "\n\n")

	evalMutex.Unlock()

	return result, nil
}
