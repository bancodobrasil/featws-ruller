package v1

import (
	goauthgin "github.com/bancodobrasil/goauth-gin"
	"github.com/gin-gonic/gin"
)

// Router ...
func Router(router *gin.RouterGroup) {
	router.Use(goauthgin.Authenticate())
	evalRouter(router.Group("/eval"))
}
