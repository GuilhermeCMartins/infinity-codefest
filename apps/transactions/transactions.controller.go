package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupTransactionsRoutes(router *gin.Engine) {
	t := router.Group("/transactions")
	{
		t.GET("/", func(c *gin.Context) {
			transactions, count, err := FindAllTransactions()
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
