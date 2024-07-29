package api

import (
	"log"
	"myapp/db"
	"myapp/ports/consumer"
	"myapp/routes"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Start() {
	r := gin.Default()

	logrus.Info("Starting server...")
	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(r)

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
