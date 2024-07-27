package consumer

import (
	"fmt"
	"log"
)

type ConsumerManager struct {
	consumers []*Consumer
}

// NewConsumerManager cria uma nova instância do ConsumerManager
func NewConsumerManager() *ConsumerManager {
	return &ConsumerManager{}
}

// AddConsumer adiciona um novo consumidor ao ConsumerManager
func (cm *ConsumerManager) AddConsumer(amqpURI, exchange, exchangeType, key, tag string, queueName QueueName) error {
	consumer, err := NewConsumer(amqpURI, exchange, exchangeType, key, tag, queueName)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	cm.consumers = append(cm.consumers, consumer)
	return nil
}

// Start inicia todos os consumidores gerenciados pelo ConsumerManager
func (cm *ConsumerManager) Start() {
	for _, c := range cm.consumers {
		go func(consumer *Consumer) {
			if err := consumer.Start(); err != nil {
				log.Printf("Consumer %s failed: %v", consumer.tag, err)
			}
		}(c)
	}
}

// Stop encerra todos os consumidores gerenciados pelo ConsumerManager
func (cm *ConsumerManager) Stop() {
	for _, c := range cm.consumers {
		if err := c.Shutdown(); err != nil {
			log.Printf("Error shutting down consumer %s: %v", c.tag, err)
		}
	}
}
