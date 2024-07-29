package routes

import (
	"myapp/apps/transactions"
	"myapp/apps/user"
	middleware "myapp/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.HandleMethodNotAllowed = true

	router.NoMethod(middleware.MethodCheckHandler())
	router.NoRoute(middleware.NotFoundHandler())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ping",
		})
	})

	user.SetupUserRoutes(router)
	transactions.SetupTransactionsRoutes(router)
}
