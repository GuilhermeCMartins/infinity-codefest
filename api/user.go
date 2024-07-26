package api

import "github.com/gin-gonic/gin"

func SetupUserRoutes(router *gin.Engine) {
	user := router.Group("/users")
	{
			user.POST("/")
	}
}