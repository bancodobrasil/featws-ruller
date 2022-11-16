package v1

import (
	"github.com/gin-gonic/gin"
)

// Router ...
func Router(router *gin.RouterGroup) {
	//router.Use(middlewares.VerifyAPIKey())
	evalRouter(router.Group("/eval"))
}
