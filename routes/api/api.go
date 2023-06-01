package api

import (
	v1 "github.com/bancodobrasil/featws-ruller/routes/api/v1"
	"github.com/gin-gonic/gin"
)

// Router ...
func Router(router *gin.RouterGroup) {
	v1.Router(router.Group("/v1"))
}
