package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func mockHomeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Test Passed")
	}
}

func Router() *gin.RouterGroup {
	router := gin.New()
	router.GET("/", mockHomeHandler())
	return &router.RouterGroup
}

func TestHome(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	Router()

}
