package api

import (
	v1 "github.com/bancodobrasil/featws-ruller/routes/api/v1"
	"github.com/gin-gonic/gin"
)

<<<<<<< HEAD
// Router sets up a router for v1 of an API using the Gin framework.
=======
// Router ...
>>>>>>> cache-ruller
func Router(router *gin.RouterGroup) {
	v1.Router(router.Group("/v1"))
}
