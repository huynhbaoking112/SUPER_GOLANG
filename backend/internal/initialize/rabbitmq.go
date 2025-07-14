package initialize

import (
	"fmt"
	"go-backend-v2/global"
	"go-backend-v2/pkg/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ() {
	cfg := global.Config.RabbitMQ
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to RabbitMQ: %s", err))
	}

	global.RabbitMQConn = conn

	publisher, err := utils.NewRabbitMQPublisher(conn, cfg.IamExchange, amqp.ExchangeTopic)
	if err != nil {
		panic(fmt.Sprintf("Failed to create EventPublisher: %s", err))
	}

	global.EventTopicPublisher = publisher

	fmt.Printf("RabbitMQ Connection Initialized - Host: %s:%d, Exchange: %s\n",
		cfg.Host, cfg.Port, cfg.IamExchange)
}
