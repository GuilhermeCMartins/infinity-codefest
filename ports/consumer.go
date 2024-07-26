package ports

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/db"
	"myapp/models"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer representa um consumidor RabbitMQ
type Consumer struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	tag           string
	done          chan error
	msgCount      int
	msgCountMutex sync.Mutex
}

// NewConsumer cria um novo consumidor RabbitMQ
func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, tag string) (*Consumer, error) {
	c := &Consumer{
		tag:  tag,
		done: make(chan error),
	}

	var err error

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("sample-consumer")
	c.conn, err = amqp.DialConfig(amqpURI, config)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	// c.channel.Qos(1, 0, false)

	if _, err = c.channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	deliveries, err := c.channel.Consume(
		queueName,
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

	go c.handle(deliveries)

	return c, nil
}

// Start inicia o consumidor RabbitMQ
func (c *Consumer) Start() error {
	return <-c.done
}

// Shutdown encerra o consumidor RabbitMQ
func (c *Consumer) Shutdown() error {
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	return nil
}

// handle processa as mensagens recebidas
// Precisamos desaclopar consumer, producer e logica de verificao e criacao
func (c *Consumer) handle(deliveries <-chan amqp.Delivery) {
	defer func() {
		c.done <- nil
	}()

	db := db.Init()
	for d := range deliveries {
		var payload models.UserPayload

		err := json.Unmarshal([]byte(d.Body), &payload)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		log.Print(payload)

		producer, err := NewProducer("amqp://guest:guest@localhost:5672/", "", "direct", "users", "test-key")
		if err != nil {
			log.Fatalf("Failed to create producer: %v", err)
			defer producer.Shutdown()
		}

		result := models.HandleMessageUser(db, payload)
		if result != "" {

			// fmt.Println()
			// fmt.Print(result)
			// fmt.Println()

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
