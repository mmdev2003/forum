package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
)

func New(
	host string,
	port string,
	username string,
	password string,
) *RabbitMQ {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)

	conn, err := amqp091.Dial(url)
	if err != nil {
		panic(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &RabbitMQ{
		channel: channel,
		conn:    conn,
	}
}

type RabbitMQ struct {
	channel *amqp091.Channel
	conn    *amqp091.Connection
}

type Event struct {
	Body any `json:"body"`
}

func (r *RabbitMQ) Publish(ctx context.Context, queueName string, body any) error {
	event := Event{
		Body: body,
	}

	jsonData, err := json.Marshal(event)

	err = r.channel.PublishWithContext(ctx, "", queueName, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		Body:         jsonData,
		DeliveryMode: amqp091.Persistent,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) Subscribe(queueName string, handler func(event amqp091.Delivery) error) {
	_, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	events, err := r.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	go func() {
		for event := range events {
			slog.Info("start handle event", "event", event)
			err = handler(event)
			if err != nil {
				slog.Error("failed handle event", "err", err)
				event.Nack(false, true)
			} else {
				slog.Info("event handled", "event", event)
				event.Ack(false)
			}
		}
	}()
}
