package awss3

import (
	"github.com/gin-gonic/gin"
)

// Router ...
func Router(router *gin.RouterGroup) {
	awsS3Router(router.Group("/aws"))
}
