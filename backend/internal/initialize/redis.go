package initialize

import (
	"context"
	"fmt"
	"go-backend-v2/global"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	cfg := global.Config.Redis

	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		DialTimeout:     time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(cfg.WriteTimeout) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.ConnMaxIdleTime) * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("failed to connect to Redis: %v", err))
	}

	global.RedisClient = client

	fmt.Printf("Redis connected successfully at %s:%d (DB: %d)\n", cfg.Host, cfg.Port, cfg.DB)
}
