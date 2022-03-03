package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sync"

	"github.com/bancodobrasil/featws-ruller/service"
	"github.com/bancodobrasil/featws-ruller/types"
	"github.com/gin-gonic/gin"
)

// DefaultKnowledgeBaseName its default name of Knowledge Base
const DefaultKnowledgeBaseName = "default"

// DefaultKnowledgeBaseVersion its default version of Knowledge Base
const DefaultKnowledgeBaseVersion = "latest"

var LoadMutex sync.Mutex

func HomeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "FeatWS Works!!!")
	}

}

func EvalHandler() gin.HandlerFunc {
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

		LoadMutex.Lock()

		knowledgeBase := service.KnowledgeLibrary.GetKnowledgeBase(knowledgeBaseName, version)

		if !knowledgeBase.ContainsRuleEntry("DefaultValues") {

			err := service.LoadRemoteGRL(knowledgeBaseName, version)
			if err != nil {
				log.Printf("Erro on load: %v", err)
				c.Status(http.StatusNotFound)
				fmt.Fprint(c.Writer, "KnowledgeBase or version not founded!")
				// w.WriteHeader(http.StatusServiceUnavailable)
				// encoder := json.NewEncoder(w)
				// encoder.Encode(err)
				LoadMutex.Unlock()
				return
			}

			knowledgeBase = service.KnowledgeLibrary.GetKnowledgeBase(knowledgeBaseName, version)

			if !knowledgeBase.ContainsRuleEntry("DefaultValues") {
				c.Status(http.StatusNotFound)
				fmt.Fprint(c.Writer, "KnowledgeBase or version not founded!")
				LoadMutex.Unlock()
				return
			}
		}

		LoadMutex.Unlock()

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

		result, err := service.Eval(ctx, knowledgeBase)
		if err != nil {
			panic(err)
		}

		// log.Print("Context:\n\t", ctx.GetEntries(), "\n\n")
		// log.Print("Features:\n\t", result.GetFeatures(), "\n\n")

		c.JSON(http.StatusOK, result.GetFeatures())
		//fmt.Fprintf(w, "%v", result)
	}

}
