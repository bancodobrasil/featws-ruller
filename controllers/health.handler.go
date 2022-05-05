package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthLiveHandler ...
func HealthLiveHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Application is Live!!!")
	}

}
