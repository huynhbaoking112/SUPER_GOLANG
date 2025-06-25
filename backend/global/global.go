package global

import (
	"go-backend-v2/pkg/setting"

	"github.com/redis/go-redis/v9"
)

var (
	Config      *setting.Config
	RedisClient *redis.Client // Redis connection
)
