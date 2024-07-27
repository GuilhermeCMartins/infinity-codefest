package api

import (
	"log"
	"myapp/db"
	"myapp/ports/consumer"
	"os"
	"os/signal"
	"syscall"

	middlewares "myapp/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Start() {
	r := gin.Default()

	logrus.Info("Starting server...")

	r.NoMethod(middlewares.MethodCheckHandler())
	r.NoRoute(middlewares.NotFoundHandler())

	db.Init()

	cm := consumer.NewConsumerManager()

	errConsumer := cm.AddConsumer(
		"amqp://guest:guest@localhost:5672/",
		"test-exchange",
		"direct",
		"test-key",
		"users-consumer",
		"users",
	)
	if errConsumer != nil {
		log.Fatalf("Failed to add consumer: %v", errConsumer)
	}

	errTransactions := cm.AddConsumer(
		"amqp://guest:guest@localhost:5672/",
		"test-exchange",
		"direct",
		"test-key",
		"users-consumer",
		"transactions",
	)

	if errTransactions != nil {
		log.Fatalf("Failed to add consumer: %v", errConsumer)
	}

	cm.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		logrus.Info("Shutting down gracefully...")

		cm.Stop()

		os.Exit(0)
	}()

	if err := r.Run(":4000"); err != nil {
		logrus.Fatalf("Server failed: %v", err)
	}
}