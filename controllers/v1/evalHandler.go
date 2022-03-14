package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sync"

	"github.com/bancodobrasil/featws-ruller/services"
	"github.com/bancodobrasil/featws-ruller/types"
	"github.com/gin-gonic/gin"
)

// DefaultKnowledgeBaseName its default name of Knowledge Base
const DefaultKnowledgeBaseName = "default"

// DefaultKnowledgeBaseVersion its default version of Knowledge Base
const DefaultKnowledgeBaseVersion = "latest"

//LoadMutex ...
var loadMutex sync.Mutex

// LoadRemoteGRL ...
var LoadRemoteGRL = services.LoadRemoteGRL

//EvalHandler ...
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

		loadMutex.Lock()

		knowledgeBase := services.KnowledgeLibrary.GetKnowledgeBase(knowledgeBaseName, version)

		if !knowledgeBase.ContainsRuleEntry("DefaultValues") {

			err := LoadRemoteGRL(knowledgeBaseName, version)
			if err != nil {
				log.Printf("Erro on load: %v", err)
				c.Status(http.StatusNotFound)
				fmt.Fprint(c.Writer, "KnowledgeBase or version not founded!")
				// w.WriteHeader(http.StatusservicesUnavailable)
				// encoder := json.NewEncoder(w)
				// encoder.Encode(err)
				loadMutex.Unlock()
				return
			}

			knowledgeBase = services.KnowledgeLibrary.GetKnowledgeBase(knowledgeBaseName, version)

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

		result, err := services.Eval(ctx, knowledgeBase)
		if err != nil {
			panic(err)
		}

		// log.Print("Context:\n\t", ctx.GetEntries(), "\n\n")
		// log.Print("Features:\n\t", result.GetFeatures(), "\n\n")

		c.JSON(http.StatusOK, result.GetFeatures())
		//fmt.Fprintf(w, "%v", result)
	}

}
