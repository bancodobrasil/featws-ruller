package routes

import (
	"github.com/bancodobrasil/featws-ruller/controllers"
	"github.com/gin-gonic/gin"
)

// HealthRouter ...
func HealthRouter(router *gin.RouterGroup) {
	router.GET("/live", controllers.HealthLiveHandler())
}
