package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MethodCheck(allowedMethods ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		methodAllowed := false
		for _, method := range allowedMethods {
			if c.Request.Method == method {
				methodAllowed = true
				break
			}
		}
		if !methodAllowed {
			c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
			c.Abort()
			return
		}
		c.Next()
	}
}
