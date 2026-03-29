package repository

import (
	"context"
	"encoding/json"
	"log"
	"notification_consumer/internal/entity"
	"notification_consumer/pkg/notification"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQRepository struct {
	conn         *amqp.Connection
	exchangeName string
}

func NewRabbitMQRepository(conn *amqp.Connection, exchange string) *RabbitMQRepository {
	return &RabbitMQRepository{conn: conn, exchangeName: exchange}
}

func (r *RabbitMQRepository) ConsumeMessagesFcm(ctx context.Context, serviceName string) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(r.exchangeName, "fanout", true, false, false, false, nil)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(serviceName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind(q.Name, "", r.exchangeName, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, serviceName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var message entity.Message
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("[%s] Failed to unmarshal message: %v", serviceName, err)
				continue
			}
			log.Printf("[%s] Received message: %+v", serviceName, message)

			notification.SendFcm(message)

			msg.Ack(false)
		}
	}()

	return nil
}

func (r *RabbitMQRepository) ConsumeMessagesSms(ctx context.Context, serviceName string) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(r.exchangeName, "fanout", true, false, false, false, nil)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(serviceName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind(q.Name, "", r.exchangeName, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, serviceName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var message entity.Message
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("[%s] Failed to unmarshal message: %v", serviceName, err)
				continue
			}
			log.Printf("[%s] Received message: %+v", serviceName, message)
			notification.SendSms(message)
			msg.Ack(false)
		}
	}()

	return nil
}

func (r *RabbitMQRepository) ConsumeMessagesEmail(ctx context.Context, serviceName string) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(r.exchangeName, "fanout", true, false, false, false, nil)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(serviceName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind(q.Name, "", r.exchangeName, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, serviceName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var message entity.Message
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("[%s] Failed to unmarshal message: %v", serviceName, err)
				continue
			}
			log.Printf("[%s] Received message: %+v", serviceName, message)

			notification.SendEmail(message)

			msg.Ack(false)
		}
	}()

	return nil
}
