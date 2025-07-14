package global

import (
	"go-backend-v2/pkg/setting"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type EventPublisher interface {
	Publish(topic string, payload interface{}) error
	Close() error
}

var (
	Config              *setting.Config
	RedisClient         *redis.Client    // Redis connection
	DB                  *gorm.DB         // MySQL database connection
	RabbitMQConn        *amqp.Connection // RabbitMQ connection
	EventTopicPublisher EventPublisher   // Event publisher service
)
