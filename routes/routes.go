package routes

import (
	"myapp/apps/transactions"
	"myapp/apps/user"
	"myapp/db"
	middlewares "myapp/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	//TODO: Fix middlewares
	router.NoMethod(middlewares.MethodCheckHandler())
	router.NoRoute(middlewares.NotFoundHandler())

	db := db.Init()

	user.SetupUserRoutes(router, db)
	transactions.SetupTransactionsRoutes(router, db)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ping")
})
}
