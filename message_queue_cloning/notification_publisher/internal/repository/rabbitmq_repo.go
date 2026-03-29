package repository

import (
	"context"
	"encoding/json"
	"notification_publisher/internal/entity"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQRepository struct {
	conn         *amqp.Connection
	exchangeName string
}

func NewRabbitMQRepository(conn *amqp.Connection, exchange string) *RabbitMQRepository {
	return &RabbitMQRepository{conn: conn, exchangeName: exchange}
}

func (r *RabbitMQRepository) PublishMessage(ctx context.Context, routingKey string, message entity.Message) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		r.exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(ctx, r.exchangeName, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
