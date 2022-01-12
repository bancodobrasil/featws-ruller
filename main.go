package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"

	"github.com/bancodobrasil/featws-ruller/processor"
	"github.com/bancodobrasil/featws-ruller/types"
	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "FeatWS Works!!!")
}

func EvalHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t map[string]interface{}
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)

	knowledgeBase := knowledgeLibrary.GetKnowledgeBase("TutorialRules", "0.0.1")

	ctx := types.NewContext()

	keys := reflect.ValueOf(t).MapKeys()

	for i := range keys {
		k := keys[i]
		ctx.Put(k.String(), t[k.String()])
	}
	result, err := eval(ctx, knowledgeBase)
	if err != nil {
		panic(err)
	}

	fmt.Print("Context:\n\t", ctx.GetEntries(), "\n\n")
	fmt.Print("Features:\n\t", result.GetFeatures(), "\n\n")

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(result.GetFeatures())
	//fmt.Fprintf(w, "%v", result)
}

var knowledgeLibrary *ast.KnowledgeLibrary = ast.NewKnowledgeLibrary()

// Hello returns a greeting for the named person.
func main() {

	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	fileRes := pkg.NewFileResource("../featws-compiler/examples/full/rules.grl")
	err := ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", fileRes)
	if err != nil {
		panic(err)
	}

	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("TutorialRules", "0.0.1")

	ctx := types.NewContext()

	ctx.Put("idade", "45")
	ctx.Put("branch", "03411")
	ctx.Put("account", "00000170408")

	_, err = eval(ctx, knowledgeBase)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/eval/", EvalHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
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

	fmt.Print("Context:\n\t", ctx.GetEntries(), "\n\n")
	fmt.Print("Features:\n\t", result.GetFeatures(), "\n\n")

	return result, nil
}
