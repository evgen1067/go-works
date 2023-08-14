package consumer

import (
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/streadway/amqp"
)

type Consumer struct {
	*rabbit.RMQ
}

func NewConsumer(uri, queue string) *Consumer {
	return &Consumer{rabbit.InitRabbitMQ(uri, queue)}
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	return c.RMQ.Chan.Consume(
		c.RMQ.Queue, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
}
