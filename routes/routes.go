package routes

import (
	"github.com/bancodobrasil/featws-ruller/config"
	"github.com/bancodobrasil/featws-ruller/docs"
	"github.com/bancodobrasil/featws-ruller/routes/api"
	"github.com/bancodobrasil/featws-ruller/routes/health"
	telemetry "github.com/bancodobrasil/gin-telemetry"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// SetupRoutes ...
func SetupRoutes(router *gin.Engine) {
	cfg := config.GetConfig()
	docs.SwaggerInfo.Host = cfg.ExternalHost
	homeRouter(router.Group("/"))
	health.Router(router.Group("/health"))
	// setup swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// APIRoutes define all api routes
func APIRoutes(router *gin.Engine) {
	group := router.Group("/api")
	group.Use(telemetry.Middleware("featws-ruller"))
	api.Router(group)
}
