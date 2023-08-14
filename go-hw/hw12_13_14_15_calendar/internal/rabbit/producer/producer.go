package producer

import (
	"context"
	"fmt"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/streadway/amqp"
)

type Producer struct {
	*rabbit.RMQ
}

func NewProducer(uri, queue string) *Producer {
	return &Producer{rabbit.InitRabbitMQ(uri, queue)}
}

func (p *Producer) Publish(ctx context.Context, body []byte) (err error) {
	err = p.RMQ.Chan.Publish(
		"",          // exchange
		p.RMQ.Queue, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message. Error: %w", err)
	}
	return nil
}
