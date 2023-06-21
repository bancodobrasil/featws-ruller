package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler returns a Gin HTTP handler function that responds with a message if
// the application is available (200 OK).
func HomeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "FeatWS Ruller Works!!!")
	}

}
