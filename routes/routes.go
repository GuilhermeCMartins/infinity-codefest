package api

import (
	user "myapp/api"
	middlewares "myapp/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.NoMethod(middlewares.MethodCheckHandler())
	router.NoRoute(middlewares.NotFoundHandler())

	user.SetupUserRoutes(router)
}