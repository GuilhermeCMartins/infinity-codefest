package api

import "github.com/gin-gonic/gin"

func SetupUserRoutes(router *gin.Engine) {
	user := router.Group("/users")
	{
		user.POST("/")
	}

	// producer, err := ports.NewProducer("amqp://guest:guest@localhost:5672/", "test-exchange", "direct", "users", "test-key")
	// if err != nil {
	// 	log.Fatalf("Failed to create producer: %v", err)
	// }
	// defer producer.Shutdown()

	// if err := producer.Publish("Hello, RabbitMQ!"); err != nil {
	// 	log.Fatalf("Failed to publish message: %v", err)
	// }
}
