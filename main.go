package main

import (
	"github.com/gin-gonic/gin"
)

// Class that start and stop server
// Initiated db

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "ping, pong")
	})

	r.GET("/users", func(c *gin.Context) {
		c.String(200, "Usuários")
	})

	r.NoRoute(func(c *gin.Context) {
		c.String(404, "Teste")
	})

	r.Run(":8080")
}
