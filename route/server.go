package route

import (
	"github.com/bancodobrasil/featws-ruller/controller"
	"github.com/bancodobrasil/featws-ruller/service"
	"github.com/gin-gonic/gin"
)

func SetupServer(router *gin.Engine) {

	router.GET("/", controller.HomeHandler())
	router.POST("/eval/:knowledgeBase/:version", controller.EvalHandler())
	router.POST("/eval/:knowledgeBase/:version/", controller.EvalHandler())
	router.POST("/eval/:knowledgeBase", controller.EvalHandler())
	router.POST("/eval/:knowledgeBase/", controller.EvalHandler())

	knowledgeBase := service.KnowledgeLibrary.GetKnowledgeBase(controller.DefaultKnowledgeBaseName, controller.DefaultKnowledgeBaseVersion)

	if knowledgeBase.ContainsRuleEntry("DefaultValues") {

		router.POST("/eval/", controller.EvalHandler())
		router.POST("/eval", controller.EvalHandler())

	}

}
