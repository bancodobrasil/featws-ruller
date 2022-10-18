package middlewares

import "github.com/gin-gonic/gin"

// Middleware ...
type Middleware interface {
	Run()
}

// InitializeMiddlewares ...
func InitializeMiddlewares() {
	NewVerifyAPIKeyMiddleware()
}

// Helper function to abort the request with an error status code and message
func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}
