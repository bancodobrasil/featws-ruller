package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/bancodobrasil/featws-ruller/types"
	"github.com/gin-gonic/gin"
)

var loadMutex sync.Mutex

func setupServer(router *gin.Engine) *http.Server {

	router.GET("/", homeHandler())
	router.GET("/eval/:knowledgeBase/:version", evalHandler())
	router.GET("/eval/:knowledgeBase/:version/", evalHandler())
	router.GET("/eval/:knowledgeBase", evalHandler())
	router.GET("/eval/:knowledgeBase/", evalHandler())

	knowledgeBase := knowledgeLibrary.GetKnowledgeBase(DefaultKnowledgeBaseName, DefaultKnowledgeBaseVersion)

	if knowledgeBase.ContainsRuleEntry("DefaultValues") {

		router.GET("/eval/", evalHandler())
		router.GET("/eval", evalHandler())

	}

	port := getEnv("PORT", "8000")

	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: router,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listen on http://0.0.0.0:%s", port)

	return srv
}

func homeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "FeatWS Works!!!")
		fmt.Fprintf(c.Writer, "FeatWS Works!!!")
	}

}

func evalHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		knowledgeBaseName := c.Param("knowledgeBase")
		if knowledgeBaseName == "" {
			knowledgeBaseName = DefaultKnowledgeBaseName
		}

		version := c.Param("version")
		if version == "" {
			version = DefaultKnowledgeBaseVersion
		}

		log.Printf("Eval with %s %s\n", knowledgeBaseName, version)

		loadMutex.Lock()

		knowledgeBase := knowledgeLibrary.GetKnowledgeBase(knowledgeBaseName, version)

		if !knowledgeBase.ContainsRuleEntry("DefaultValues") {

			err := loadRemoteGRL(knowledgeBaseName, version)
			if err != nil {
				log.Printf("Erro on load: %v", err)
				c.Status(http.StatusNotFound)
				fmt.Fprint(c.Writer, "KnowledgeBase or version not founded!")
				// w.WriteHeader(http.StatusServiceUnavailable)
				// encoder := json.NewEncoder(w)
				// encoder.Encode(err)
				loadMutex.Unlock()
				return
			}

			knowledgeBase = knowledgeLibrary.GetKnowledgeBase(knowledgeBaseName, version)

			if !knowledgeBase.ContainsRuleEntry("DefaultValues") {
				c.Status(http.StatusNotFound)
				fmt.Fprint(c.Writer, "KnowledgeBase or version not founded!")
				loadMutex.Unlock()
				return
			}
		}

		loadMutex.Unlock()

		decoder := json.NewDecoder(c.Request.Body)
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

		c.Status(http.StatusOK)
		encoder := json.NewEncoder(c.Writer)
		encoder.Encode(result.GetFeatures())
		//fmt.Fprintf(w, "%v", result)
	}

}
