package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	r.GET("/about", func(c *gin.Context) {
		c.String(200, "About Page")
	})

	r.Run(":8080")
}
