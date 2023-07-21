package v1

import (
	goauthgin "github.com/bancodobrasil/goauth-gin"
	"github.com/gin-gonic/gin"
)

// Router sets up a router with authentication middleware and a sub-router for evaluating code.
func Router(router *gin.RouterGroup) {
	router.Use(goauthgin.Authenticate())
	evalRouter(router.Group("/eval"))
}
