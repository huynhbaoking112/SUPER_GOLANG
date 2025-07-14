package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"go-backend-v2/internal/common"
	"go-backend-v2/internal/dto"
	"time"

	"go-backend-v2/global"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func NewRabbitMQPublisher(conn *amqp.Connection, exchangeName string, exchangeType string) (global.EventPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &rabbitMQPublisher{
		conn:     conn,
		channel:  ch,
		exchange: exchangeName,
	}, nil
}

func (p *rabbitMQPublisher) Publish(topic string, payload interface{}) error {
	event := dto.GenericEvent{
		EventID:       uuid.New().String(),
		Topic:         topic,
		SourceService: common.SourceServiceIAM,
		Timestamp:     time.Now().UTC(),
		Payload:       payload,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.channel.PublishWithContext(context.Background(),
		p.exchange,
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (p *rabbitMQPublisher) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
	}
	return nil
}
