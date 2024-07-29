package api

import (
	"log"
	"myapp/ports/consumer"
	"myapp/routes"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Start() {
	r := gin.Default()

	logrus.Info("Starting server...")

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))

	routes.SetupRoutes(r)

	userSettings := consumer.ConsumerSettings{
		AMQPURI:      "amqp://guest:guest@localhost:5672/",
		Exchange:     "test-exchange",
		ExchangeType: "direct",
		Key:          "test-key",
		Tag:          "users-consumer",
		QueueName:    consumer.QueueName("users"),
	}

	transactionSettings := consumer.ConsumerSettings{
		AMQPURI:      "amqp://guest:guest@localhost:5672/",
		Exchange:     "test-exchange",
		ExchangeType: "direct",
		Key:          "test-key",
		Tag:          "transactions-consumer",
		QueueName:    consumer.QueueName("transactions"),
	}

	cmUsers, err := consumer.NewConsumerManager(1, userSettings)
	if err != nil {
		log.Fatalf("Error creating Users ConsumerManager: %v", err)
	}

	cmTransaction, err := consumer.NewConsumerManager(1, transactionSettings)
	if err != nil {
		log.Fatalf("Error creating Transactions ConsumerManager: %v", err)
	}

	cmUsers.Start()
	cmTransaction.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		logrus.Info("Shutting down gracefully...")

		cmUsers.Stop()
		cmTransaction.Stop()

		os.Exit(0)
	}()

	if err := r.Run(":4000"); err != nil {
		logrus.Fatalf("Server failed: %v", err)
	}
}
