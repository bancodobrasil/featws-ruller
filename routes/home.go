package routes

import (
	"github.com/bancodobrasil/featws-ruller/controllers"
	"github.com/gin-gonic/gin"
)

// homeRouter sets up a route for the home page using the GET method and calls the HomeHandler
// function from the controllers package.
func homeRouter(router *gin.RouterGroup) {
	router.GET("/", controllers.HomeHandler())
}
