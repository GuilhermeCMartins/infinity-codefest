package consumer

import (
	"fmt"
	"log"
)

type ConsumerSettings struct {
	AMQPURI      string
	Exchange     string
	ExchangeType string
	Key          string
	Tag          string
	QueueName    QueueName
}
type ConsumerManager struct {
	consumers []*Consumer
}

func NewConsumerManager(numConsumers int, config ConsumerSettings) (*ConsumerManager, error) {
	cm := &ConsumerManager{}

	for i := 0; i < numConsumers; i++ {
		consumerConfig := config
		consumer, err := NewConsumer(consumerConfig.AMQPURI, consumerConfig.Exchange, consumerConfig.ExchangeType, consumerConfig.Key, consumerConfig.Tag, consumerConfig.QueueName)
		if err != nil {
			return nil, fmt.Errorf("failed to create consumer for queue %s: %w", consumerConfig.QueueName, err)
		}
		cm.consumers = append(cm.consumers, consumer)
	}

	return cm, nil
}

func (cm *ConsumerManager) Start() {
	for _, c := range cm.consumers {
		go func(consumer *Consumer) {
			if err := consumer.Start(); err != nil {
				log.Printf("Consumer %s failed: %v", consumer.tag, err)
			}
		}(c)
	}
}

func (cm *ConsumerManager) Stop() {
	for _, c := range cm.consumers {
		if err := c.Shutdown(); err != nil {
			log.Printf("Error shutting down consumer %s: %v", c.tag, err)
		}
	}
}
