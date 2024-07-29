package routes

import (
	"myapp/apps/user"
	"myapp/db"
	middlewares "myapp/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.NoMethod(middlewares.MethodCheckHandler())
	router.NoRoute(middlewares.NotFoundHandler())

	db := db.Init()

	user.SetupUserRoutes(router, db)


}
