package configure

import (
	"github.com/go-redis/redis"

	"github.com/b0pof/ppo/internal/config"
)

const defaultDB = 0

func MustInitRedis(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		DB:   defaultDB,
	})
}
