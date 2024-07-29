package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTransactionsRoutes(router *gin.Engine, db *gorm.DB) {
	t := router.Group("/transactions")
	{
		t.GET("/", func(c *gin.Context) {
			transactions, count, err := FindAllTransactions(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"transactions": transactions,
				"count":        count,
			})
		})

	}
}
