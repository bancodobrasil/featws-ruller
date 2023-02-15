package api

import (
	awss3 "github.com/bancodobrasil/featws-ruller/routes/api/aws-s3"
	v1 "github.com/bancodobrasil/featws-ruller/routes/api/v1"
	"github.com/gin-gonic/gin"
)

// Router ...
func Router(router *gin.RouterGroup) {
	v1.Router(router.Group("/v1"))
	awss3.Router(router.Group("/awss3"))
}
