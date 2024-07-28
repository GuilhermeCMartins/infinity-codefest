package consumer

import (
	"fmt"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

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
