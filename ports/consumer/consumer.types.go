package consumer

import (
	"sync"

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
