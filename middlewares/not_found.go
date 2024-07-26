package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	}
}
