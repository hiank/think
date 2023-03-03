package defaults

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func RedisMasterOptions() *redis.Options {
	///
	return &redis.Options{
		Addr:     "redis-master:tcp-redis",
		Password: os.Getenv("redis-password"),
		DB:       0,
	}
}

func RedisSlaveOptions() *redis.Options {
	return &redis.Options{
		Addr:     "redis-slave:tcp-redis",
		Password: os.Getenv("redis-password"),
		DB:       0,
	}
}

func NatsUrl() string {
	return "nats:tcp-nats"
}
