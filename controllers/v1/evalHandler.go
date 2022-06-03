package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	payloads "github.com/bancodobrasil/featws-ruller/payloads/v1"
	"github.com/bancodobrasil/featws-ruller/services"
	"github.com/bancodobrasil/featws-ruller/types"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//LoadMutex ...
var loadMutex sync.Mutex

// EvalHandler godoc
// @Summary 		Evaluate the rulesheet
// @Description 	Receive the params to execute the rulesheet
// @Tags 			eval
// @Accept  		json
// @Produce  		json
// @Param			knowledgeBase path string false "knowledgeBase"
// @Param 			version path string false "version"
// @Param  			parameters body payloads.Eval true "Parameters"
// @Success 		200 {string} string "ok"
// @Failure 		400,404 {object} string
// @Failure 		500 {object} string
// @Failure 		default {object} string
// @Router 			/eval/{knowledgeBase}/{version} [post]
// @Router 			/eval/{knowledgeBase} [post]
// @Router 			/eval [post]
func EvalHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		knowledgeBaseName := c.Param("knowledgeBase")
		if knowledgeBaseName == "" {
			knowledgeBaseName = services.DefaultKnowledgeBaseName
		}

		version := c.Param("version")
		if version == "" {
			version = services.DefaultKnowledgeBaseVersion
		}

		log.Debugf("Eval with %s %s\n", knowledgeBaseName, version)

		loadMutex.Lock()

		knowledgeBase := services.EvalService.GetKnowledgeLibrary().GetKnowledgeBase(knowledgeBaseName, version)

		if !(len(knowledgeBase.RuleEntries) > 0) {

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

			if !(len(knowledgeBase.RuleEntries) > 0) {
				c.Status(http.StatusNotFound)
				fmt.Fprint(c.Writer, "KnowledgeBase or version not founded!")
				loadMutex.Unlock()
				return
			}
		}

		loadMutex.Unlock()

		decoder := json.NewDecoder(c.Request.Body)
		var t payloads.Eval
		err := decoder.Decode(&t)
		if err != nil {
			log.Errorf("Erro on json decode: %v", err)
			c.Status(http.StatusInternalServerError)
			fmt.Fprint(c.Writer, "Error on json decode")
			return
		}
		log.Debugln(t)

		ctx := types.NewContextFromMap(t)

		result, err := services.EvalService.Eval(ctx, knowledgeBase)
		requiredError := errors.New("rule engine execute panic ! recovered : &{0xc000162000 map[] 2022-06-02 18:07:46.9685124 -0300 -03 m=+111.767871001 panic <nil> The param is not registry as remote loaded and is required <nil> <nil> }")
		log.Print(requiredError)
		log.Print(err)
		if err != nil {
			if errors.Is(err, requiredError) {
				log.Errorf("Error on eval: %v", err)
				c.Status(http.StatusBadRequest)
				fmt.Fprintf(c.Writer, "The param is not registry as remote loaded and is required")
				return
			}

			log.Errorf("Error on eval: %v", err)
			c.Status(http.StatusInternalServerError)
			fmt.Fprintf(c.Writer, "Error on eval: %v", err)
			return
		}

		log.Debug("Context:\n\t", ctx.GetEntries(), "\n\n")
		log.Debug("Features:\n\t", result.GetFeatures(), "\n\n")

		c.JSON(http.StatusOK, result.GetFeatures())
	}

}
