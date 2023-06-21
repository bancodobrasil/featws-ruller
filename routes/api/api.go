package api

import (
	v1 "github.com/bancodobrasil/featws-ruller/routes/api/v1"
	"github.com/gin-gonic/gin"
)

// Router sets up a router for v1 of an API using the Gin framework.
func Router(router *gin.RouterGroup) {
	v1.Router(router.Group("/v1"))
}
