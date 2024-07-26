package ports

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchange     string
	exchangeType string
	queueName    string
	routingKey   string
}

func NewProducer(amqpURI, exchange, exchangeType, queueName, routingKey string) (*Producer, error) {
	p := &Producer{
		exchange:     exchange,
		exchangeType: exchangeType,
		queueName:    queueName,
		routingKey:   routingKey,
	}

	var err error
	p.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	if _, err = p.channel.QueueDeclare(
		p.queueName,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	return p, nil
}

func (p *Producer) Publish(body string) error {
	fmt.Print(body)
	if err := p.channel.Publish(
		"",
		"users",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	); err != nil {
		return fmt.Errorf("Publish: %s", err)
	}

	// Wait for the confirmation
	// confirm := <-p.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	// if !confirm.Ack {
	// 	return fmt.Errorf("Publish failed")
	// }

	// fmt.Print(confirm.Ack)

	return nil
}

func (p *Producer) Shutdown() error {
	if err := p.channel.Close(); err != nil {
		return fmt.Errorf("Channel Close: %s", err)
	}

	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("Connection Close: %s", err)
	}

	return nil
}
