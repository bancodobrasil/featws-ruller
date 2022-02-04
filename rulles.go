package main

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"

	"github.com/bancodobrasil/featws-ruller/processor"
	"github.com/bancodobrasil/featws-ruller/types"
)

var knowledgeLibrary *ast.KnowledgeLibrary = ast.NewKnowledgeLibrary()

func loadLocalGRL(grlPath string, knowledgeBaseName string, version string) error {
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	fileRes := pkg.NewFileResource(grlPath)
	return ruleBuilder.BuildRuleFromResource(knowledgeBaseName, version, fileRes)
}

type KnowledgeBaseInfo struct {
	KnowledgeBaseName string
	Version           string
}

func loadRemoteGRL(knowledgeBaseName string, version string) error {
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	headers := make(http.Header)
	for header, value := range config.ResourceLoader.Headers {
		headers.Set(header, value)
	}

	url := config.ResourceLoader.Url
	url = strings.Replace(url, "{knowledgeBase}", "{{.KnowledgeBaseName}}", -1)
	url = strings.Replace(url, "{version}", "{{.Version}}", -1)

	info := KnowledgeBaseInfo{
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

func eval(ctx *types.Context, knowledgeBase *ast.KnowledgeBase) (*types.Result, error) {
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

	return result, nil
}
