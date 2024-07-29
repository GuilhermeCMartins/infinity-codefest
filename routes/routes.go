package routes

import (
	"myapp/apps/transactions"
	"myapp/apps/user"
	"myapp/db"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	db := db.Init()

	user.SetupUserRoutes(router, db)
	transactions.SetupTransactionsRoutes(router, db)
}
