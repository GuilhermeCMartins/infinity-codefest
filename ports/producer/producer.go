package producer

import (
	"fmt"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchange     string
	exchangeType string
	queueName    string
	routingKey   string
	clientId     string
}

func NewProducer(amqpURI, exchange, exchangeType, queueName, routingKey string) (*Producer, error) {
	p := &Producer{
		exchange:     exchange,
		exchangeType: exchangeType,
		queueName:    queueName,
		routingKey:   routingKey,
		clientId:     uuid.NewString(),
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

	if err = p.channel.ExchangeDeclare(
		queueName,           // Name of the exchange
		amqp.ExchangeFanout, // Type of the exchange
		false,               // Durable
		false,               // Auto-deleted
		false,               // Internal
		false,               // No-wait
		nil,                 // Arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	return p, nil
}

func (p *Producer) Publish(body string) error {
	messageId := fmt.Sprintf("hash-ng-%s", uuid.NewString())

	if err := p.channel.Publish(
		p.queueName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			MessageId:   messageId,
		},
	); err != nil {
		return fmt.Errorf("Publish: %s", err)
	}

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
