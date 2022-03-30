package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/bancodobrasil/featws-ruller/services"
	"github.com/bancodobrasil/featws-ruller/types"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// DefaultKnowledgeBaseName its default name of Knowledge Base
const DefaultKnowledgeBaseName = "default"

// DefaultKnowledgeBaseVersion its default version of Knowledge Base
const DefaultKnowledgeBaseVersion = "latest"

//LoadMutex ...
var loadMutex sync.Mutex

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

		knowledgeBase := services.EvalService.GetKnowledgeLibrary().GetKnowledgeBase(knowledgeBaseName, version)

		if !knowledgeBase.ContainsRuleEntry("DefaultValues") {

			err := services.EvalService.LoadRemoteGRL(knowledgeBaseName, version)
			if err != nil {
				log.Errorf("Erro on load: %v", err)
				c.String(http.StatusInternalServerError, "Error on load knowledgeBase and/or version")
				// w.WriteHeader(http.StatusservicesUnavailable)
				// encoder := json.NewEncoder(w)
				// encoder.Encode(err)
				loadMutex.Unlock()
				return
			}

			knowledgeBase = services.EvalService.GetKnowledgeLibrary().GetKnowledgeBase(knowledgeBaseName, version)

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
			log.Errorf("Erro on json decode: %v", err)
			c.Status(http.StatusInternalServerError)
			fmt.Fprint(c.Writer, "Error on json decode")
			return
		}
		log.Println(t)

		ctx := types.NewContextFromMap(t)

		result, err := services.EvalService.Eval(ctx, knowledgeBase)
		if err != nil {
			log.Errorf("Error on eval: %v", err)
			c.Status(http.StatusInternalServerError)
			fmt.Fprint(c.Writer, "Error on eval")
			return
		}

		// log.Print("Context:\n\t", ctx.GetEntries(), "\n\n")
		// log.Print("Features:\n\t", result.GetFeatures(), "\n\n")

		c.JSON(http.StatusOK, result.GetFeatures())
		//fmt.Fprintf(w, "%v", result)
	}

}
