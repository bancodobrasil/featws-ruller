package v1

import (
	v1 "github.com/bancodobrasil/featws-ruller/controllers/v1"
	"github.com/bancodobrasil/featws-ruller/services"
	"github.com/gin-gonic/gin"
)

func evalRouter(router *gin.RouterGroup) {
	router.POST("/:knowledgeBase/:version", v1.EvalHandler())
	router.POST("/:knowledgeBase/:version/", v1.EvalHandler())
	router.POST("/:knowledgeBase", v1.EvalHandler())
	router.POST("/:knowledgeBase/", v1.EvalHandler())

	knowledgeBase := services.EvalService.GetDefaultKnowledgeBase()

	if len(knowledgeBase.RuleEntries) > 0 {

		router.POST("/", v1.EvalHandler())
		router.POST("", v1.EvalHandler())

	}

}
