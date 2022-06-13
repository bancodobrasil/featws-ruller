package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//HomeHandler ...
func HomeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "FeatWS Ruller Works!!!")
	}

}
