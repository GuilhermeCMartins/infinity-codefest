package middleware

import (
	"myapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func MethodCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		logrus.Warnf("%s %s %d %s", c.Request.Method, c.Request.URL, http.StatusMethodNotAllowed, utils.ErrMethodNotAllowed.Error())
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": utils.ErrMethodNotAllowed.Error()})
	}
}
