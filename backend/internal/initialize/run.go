package initialize

import (
	"fmt"
)

// Run initializes all components in the correct order
func Run() {
	fmt.Println("Initializing application...")

	// Load configuration first
	LoadConfig()
	fmt.Println("Configuration loaded")

	// Initialize database connection
	InitMysql()

	// Run database migrations
	InitMigrations()

	// Initialize Redis connection
	InitRedis()

	// Initialize RabbitMQ connection
	InitRabbitMQ()

	// Initialize logger (if implemented)
	// InitLogger()

	fmt.Println("All components initialized successfully!")
}
