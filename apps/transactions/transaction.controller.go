package transaction

import (
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {
	t := router.Group("/transactions")
	{
		t.POST("/")
	}
}
