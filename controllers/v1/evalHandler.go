package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	payloads "github.com/bancodobrasil/featws-ruller/payloads/v1"
	"github.com/bancodobrasil/featws-ruller/services"
	"github.com/bancodobrasil/featws-ruller/types"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// LoadMutex ...
var loadMutex sync.Mutex

// EvalHandler godoc
// @Summary 		Evaluate the rulesheet / Avaliação da folha de Regra
// @Description 	Ao receber os parâmetros para executar as folhas de regras. A seguir é explicado com mais detalhes sobre os endpoints:
// @Description
// @Description  	- /Eval: Esse endpoint é utilizado apenas para aplicações que possuem uma única folha de regra padrão.
// @Description  	- /Eval/{knowledgeBase}: Nesse endpoint é necessário colocar o parametro do nome da folha de regra
// @Description  	- /Eval/{knowledgeBase}/{version}: Nesse endpoint é necessário colocar o parametro do nome da folha de regra como também o número da versão da folha de regra que você deseja escrever a regra.
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
// @Security 		Authentication Api Key
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
		ctx.RawContext = c.Request.Context()

		result, err := services.EvalService.Eval(ctx, knowledgeBase)
		if err != nil {

			log.Errorf("Error on eval: %v", err)
			c.Status(http.StatusInternalServerError)
			fmt.Fprint(c.Writer, "Error on eval")
			return
		}

		log.Debug("Context:\n\t", ctx.GetEntries(), "\n\n")
		log.Debug("Features:\n\t", result.GetFeatures(), "\n\n")

		responseCode := http.StatusOK

		if result.Has("requiredParamErrors") {
			responseCode = http.StatusBadRequest
		}

		c.JSON(responseCode, result.GetFeatures())
	}

}
