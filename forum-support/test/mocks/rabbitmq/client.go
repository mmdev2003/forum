package rabbitmq

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
)

func New() *RabbitMQ {
	return &RabbitMQ{}
}

type RabbitMQ struct{}

type Event struct {
	Body any `json:"body"`
}

func (r *RabbitMQ) Publish(ctx context.Context, queueName string, body any) error {
	return nil
}

func (r *RabbitMQ) Subscribe(queueName string, handler func(event amqp091.Delivery) error) {}
