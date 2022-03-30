package routes

import (
	"github.com/bancodobrasil/featws-ruller/routes/api"
	"github.com/gin-gonic/gin"
)

//SetupRoutes ...
func SetupRoutes(router *gin.Engine) {
	homeRouter(router.Group("/"))
	api.Router(router.Group("/api"))
}
