package v1

import "github.com/gin-gonic/gin"

//Router ...
func Router(router *gin.RouterGroup) {
	evalRouter(router.Group("/eval"))
}
