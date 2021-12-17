package main

import (
	"fmt"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"

	"github.com/bancodobrasil/featws-ruller/processor"
	"github.com/bancodobrasil/featws-ruller/types"
)

// Hello returns a greeting for the named person.
func main() {

	//log.SetLevel(log.TraceLevel)

	dataCtx := ast.NewDataContext()

	processor := processor.NewProcessor()

	err := dataCtx.Add("processor", processor)
	if err != nil {
		panic(err)
	}

	ctx := types.NewContext()

	ctx.Put("idade", "55")

	err = dataCtx.Add("ctx", ctx)
	if err != nil {
		panic(err)
	}

	result := types.NewResult()

	err = dataCtx.Add("result", result)
	if err != nil {
		panic(err)
	}

	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	fileRes := pkg.NewFileResource("../featws-compiler/rules.grl")
	err = ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", fileRes)
	if err != nil {
		panic(err)
	}

	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("TutorialRules", "0.0.1")

	eng := engine.NewGruleEngine()
	err = eng.Execute(dataCtx, knowledgeBase)
	if err != nil {
		panic(err)
	}

	fmt.Print("Context:\n\t", ctx.GetEntries(), "\n\n")
	fmt.Print("Features:\n\t", result.GetFeatures(), "\n\n")
}
