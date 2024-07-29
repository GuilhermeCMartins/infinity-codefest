package user

import (
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
	
		// u.GET("/users/:id", FindUserById(db))
	}
}
