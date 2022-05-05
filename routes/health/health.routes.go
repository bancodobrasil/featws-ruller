package health

import (
	"github.com/bancodobrasil/featws-ruller/controllers"
	"github.com/gin-gonic/gin"
)

func HealthRouter(router *gin.RouterGroup) {
	router.GET("/live", controllers.HealthLiveHandler())
}
