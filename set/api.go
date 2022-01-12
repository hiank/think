package set

import (
	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/config"
	"github.com/hiank/think/set/db"
	"github.com/nats-io/nats.go"
)

type IOpenApi interface {
	//db-redis
	RedisCli(db.RedisTag) (*redis.Client, bool)

	//config-parser
	ConfigParser() config.IParser

	//message queue
	Nats() *nats.Conn
}
