package initialize

import (
	"fmt"
	"go-backend-v2/global"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
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

	// Initialize logger (if implemented)
	// InitLogger()

	fmt.Println("All components initialized successfully!")
}

// GetDB returns the global database instance
func GetDB() *gorm.DB {
	return global.DB
}

// GetRedis returns the global Redis client
func GetRedis() *redis.Client {
	return global.RedisClient
}
