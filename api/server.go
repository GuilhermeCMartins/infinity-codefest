package api

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	middlewares "myapp/middlewares"
	"myapp/ports"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Start() {
	r := gin.Default()

	logrus.Info("Starting server...")

	r.NoMethod(middlewares.MethodCheckHandler())
	r.NoRoute(middlewares.NotFoundHandler())

	cm := ports.NewConsumerManager()

	errConsumer := cm.AddConsumer(
		"amqp://guest:guest@localhost:5672/",
		"test-exchange",
		"direct",
		"users",
		"test-key",
		"users-consumer",
	)
	if errConsumer != nil {
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
