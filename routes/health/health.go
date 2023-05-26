package health

import (
	"github.com/bancodobrasil/featws-ruller/controllers"
	"github.com/gin-gonic/gin"
)

// Router sets up two routes for health checks using the Gin framework in Go.
func Router(router *gin.RouterGroup) {

	healthController := controllers.NewHealthController()
	router.GET("/live", healthController.HealthLiveHandler())
	router.GET("/ready", healthController.HealthReadyHandler())
}
