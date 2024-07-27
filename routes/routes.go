package api

import (
	"myapp/apps/user"
	middlewares "myapp/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.NoMethod(middlewares.MethodCheckHandler())
	router.NoRoute(middlewares.NotFoundHandler())

	user.SetupUserRoutes(router)
}
