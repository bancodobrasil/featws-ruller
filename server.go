package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/bancodobrasil/featws-ruller/types"
	"github.com/gorilla/mux"
)

func setupServer() *http.Server {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/eval/{knowledgeBase}/{version}", evalHandler)
	r.HandleFunc("/eval/{knowledgeBase}/{version}/", evalHandler)
	r.HandleFunc("/eval/{knowledgeBase}", evalHandler)
	r.HandleFunc("/eval/{knowledgeBase}/", evalHandler)

	knowledgeBase := knowledgeLibrary.GetKnowledgeBase(DefaultKnowledgeBaseName, DefaultKnowledgeBaseVersion)

	if knowledgeBase.ContainsRuleEntry("DefaultValues") {

		r.HandleFunc("/eval/", evalHandler)
		r.HandleFunc("/eval", evalHandler)

	}

	port := getEnv("PORT", "8000")

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listen on http://0.0.0.0:%s", port)

	return srv
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "FeatWS Works!!!")
}

func evalHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	knowledgeBaseName, ok := vars["knowledgeBase"]
	if !ok {
		knowledgeBaseName = DefaultKnowledgeBaseName
	}

	version, ok := vars["version"]
	if !ok {
		version = DefaultKnowledgeBaseVersion
	}

	log.Printf("Eval with %s %s\n", knowledgeBaseName, version)

	knowledgeBase := knowledgeLibrary.GetKnowledgeBase(knowledgeBaseName, version)

	if !knowledgeBase.ContainsRuleEntry("DefaultValues") {

		err := loadRemoteGRL(knowledgeBaseName, version)
		if err != nil {
			log.Printf("Erro on load: %w", err)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "KnowledgeBase or version not founded!")
			// w.WriteHeader(http.StatusServiceUnavailable)
			// encoder := json.NewEncoder(w)
			// encoder.Encode(err)
			return
		}

		knowledgeBase = knowledgeLibrary.GetKnowledgeBase(knowledgeBaseName, version)

		if !knowledgeBase.ContainsRuleEntry("DefaultValues") {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "KnowledgeBase or version not founded!")
			return
		}
	}

	decoder := json.NewDecoder(req.Body)
	var t map[string]interface{}
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)

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

	// log.Print("Context:\n\t", ctx.GetEntries(), "\n\n")
	// log.Print("Features:\n\t", result.GetFeatures(), "\n\n")

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(result.GetFeatures())
	//fmt.Fprintf(w, "%v", result)
}
