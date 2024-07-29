package routes

import (
	"myapp/apps/transactions"
	"myapp/apps/user"
	"myapp/db"
	middleware "myapp/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	db := db.Init()
	router.HandleMethodNotAllowed = true

	router.NoMethod(middleware.MethodCheckHandler())
	router.NoRoute(middleware.NotFoundHandler())
	
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
				"message": "ping",
		})
	})

	user.SetupUserRoutes(router, db)
	transactions.SetupTransactionsRoutes(router, db)
}
