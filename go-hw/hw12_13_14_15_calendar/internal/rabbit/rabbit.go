package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
)

type RMQ struct {
	uri   string
	Queue string
	Conn  *amqp.Connection
	Chan  *amqp.Channel
}

func InitRabbitMQ(uri, queue string) *RMQ {
	return &RMQ{
		uri:   uri,
		Queue: queue,
	}
}

func (r *RMQ) Start() error {
	var err error
	r.Conn, err = amqp.Dial(r.uri) // Создаем подключение к RabbitMQ
	if err != nil {
		return fmt.Errorf("unable to open connect to RabbitMQ server. Error: %w", err)
	}
	r.Chan, err = r.Conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel. Error: %w", err)
	}

	_, err = r.Chan.QueueDeclare(
		r.Queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a Queue. Error: %w", err)
	}

	return nil
}

func (r *RMQ) Stop() {
	_ = r.Chan.Close() // Закрываем канал
	_ = r.Conn.Close() // Закрываем подключение
}
