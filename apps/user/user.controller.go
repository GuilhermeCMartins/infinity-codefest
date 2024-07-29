package user

import (
	"myapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


func SetupUserRoutes(router *gin.Engine, db *gorm.DB) {
	u := router.Group("/users")

	{
		{
			u.GET("/", func(c *gin.Context) {
				users, count, err := FindAllUsers(db)
				
				if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
				}
				
				c.JSON(http.StatusOK, gin.H{
						"users": users,
						"count": count,
				})
			})
		}
		{
			u.GET("/:id", func(c *gin.Context) {
				id := c.Param("id")
				userId, _ := uuid.Parse(id)

				user, err := FindUserById(db, userId)
				
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, user)
			})
		}
		{
			u.GET("/:id/transactions", func(c *gin.Context) {
				id := c.Param("id")
				userId, _ := uuid.Parse(id)
				transactions, count, err := FindUserTransactions(db, userId)
				
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{
						"user_id": userId,
						"count": count,
						"transactions": transactions,
				})
			})
		}
		{
			u.GET("/:id/transactions/:tx", func(c *gin.Context) {
				id := c.Param("id")
				txId := c.Param("tx")
				userId, _ := uuid.Parse(id)
				txUUID, _ := uuid.Parse(txId)
				transaction, sender, err := FindUserTransactionByTransactionId(db, userId, txUUID)

				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusNotFound, gin.H{"message": "User or transaction not found"})
					return
				}
				
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				} 

				c.JSON(http.StatusOK, gin.H{
						"user_id": userId,
						"id": txUUID,
						"sender": sender,
						"transaction": transaction,
				})
			})
		}
		{
			u.GET("/:id/transactions/status/:status", func(c *gin.Context) {
				id := c.Param("id")
				status := c.Param("status")
				userId, _ := uuid.Parse(id)
				transactionStatus := models.TransactionStatus(status)
				transactions, count, err := FindUserTransactionsByStatus(db, userId, transactionStatus)

				//TO-DO: Refactor this to use a switch statement
				if (status != "approved" && status != "success" && status != "failed" && status != "review") {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
					return
				}
				
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusNotFound, gin.H{"message": "Invalid status"})
					return
				}

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{
						"user_id": userId,
						"status": status,
						"count": count,
						"transactions": transactions,
				})
			})
		}
	}
}
