package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/apps/user"
	"myapp/db"
	"myapp/ports/producer"
	"strings"
	"sync"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	tag           string
	done          chan error
	msgCount      int
	msgCountMutex sync.Mutex
	clientId      string
}

type QueueName string

const (
	USERS        QueueName = "users"
	TRANSACTIONS QueueName = "transactions"
)

func NewConsumer(amqpURI, exchange, exchangeType, key, tag string, queueName QueueName) (*Consumer, error) {
	c := &Consumer{
		done:     make(chan error),
		clientId: uuid.NewString(),
	}

	var err error

	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	q, err := c.channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	err = c.channel.QueueBind(
		q.Name,
		"",
		string(queueName),
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	deliveries, err := c.channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	switch queueName {
	case TRANSACTIONS:
		go c.handleTransactions(deliveries)
	case USERS:
		go c.handleUsers(deliveries)
	}

	return c, nil
}

func (c *Consumer) handleTransactions(deliveries <-chan amqp.Delivery) {

	for d := range deliveries {

		if strings.HasPrefix(d.MessageId, "hash-ng-") {
			continue
		}

		fmt.Println("Transactions: ", string(d.Body))

		if err := d.Ack(false); err != nil {
			log.Printf("Failed to acknowledge message: %v", err)
		}
	}
}

// handle processa as mensagens recebidas
// Precisamos desaclopar consumer, producer e logica de verificao e criacao
func (c *Consumer) handleUsers(deliveries <-chan amqp.Delivery) {

	db := db.Init()
	for d := range deliveries {
		var payload user.UserPayload

		if strings.HasPrefix(d.MessageId, "hash-ng-") {
			// fmt.Println("Skipping own message:", d.MessageId)
			continue
		}

		// fmt.Println("Receive new message: ", string(d.Body))

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

		result := user.HandleMessageUser(db, payload)
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

func (c *Consumer) Start() error {
	return <-c.done
}

func (c *Consumer) Shutdown() error {
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	return nil
}
