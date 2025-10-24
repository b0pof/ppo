package configure

import (
	"github.com/go-redis/redis"

	"git.iu7.bmstu.ru/kia22u475/ppo/internal/config"
)

const defaultDB = 0

func MustInitRedis(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		DB:   defaultDB,
	})
}
