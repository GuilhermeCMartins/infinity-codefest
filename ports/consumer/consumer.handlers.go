package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/apps/transactions"
	"myapp/apps/user"
	"myapp/models"
	"myapp/ports/producer"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Consumer) handleTransactions(deliveries <-chan amqp.Delivery) {

	for d := range deliveries {
		var payload models.TransactionPayload

		if strings.HasPrefix(d.MessageId, "hash-ng-") {
			fmt.Println("[TRANSACTIONS] Skipping own message:", d.MessageId)
			continue
		}

		fmt.Println("[TRANSACTIONS] Receive new message: ", string(d.Body))

		err := json.Unmarshal(d.Body, &payload)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		producer, err := producer.NewProducer("amqp://guest:guest@localhost:5672/", "", "direct", "transactions", "test-key")
		if err != nil {
			log.Fatalf("Failed to create producer: %v", err)
			defer producer.Shutdown()
		}

		result := transactions.HandleMessageTransaction(payload)
		if result != "" {

			err := producer.Publish(result)

			if err != nil {
				log.Fatalf("Failed to publish message: %v", err)
			}

		}

		if err := d.Ack(false); err != nil {
			log.Printf("Failed to acknowledge message: %v", err)
		}
	}
}

func (c *Consumer) handleUsers(deliveries <-chan amqp.Delivery) {

	for d := range deliveries {
		var payload models.UserPayload

		if strings.HasPrefix(d.MessageId, "hash-ng-") {
			fmt.Println("[USERS] Skipping own message:", d.MessageId)
			continue
		}

		fmt.Println("[USERS] Receive new message: ", string(d.Body))

		err := json.Unmarshal(d.Body, &payload)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		producer, err := producer.NewProducer("amqp://guest:guest@localhost:5672/", "", "direct", "users", "test-key")
		if err != nil {
			log.Fatalf("Failed to create producer: %v", err)
			defer producer.Shutdown()
		}

		result := user.HandleMessageUser(payload)
		if result != "" {

			err := producer.Publish(result)

			if err != nil {
				log.Fatalf("Failed to publish message: %v", err)
			}

		}

		if err := d.Ack(false); err != nil {
			log.Printf("Failed to acknowledge message: %v", err)
		}
	}
}
