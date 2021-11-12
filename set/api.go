package set

import (
	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/config"
	"github.com/nats-io/nats.go"
)

type IOpenApi interface {
	//db-redis
	RedisMasterCli() (*redis.Client, bool)
	RedisSlaveCli() (*redis.Client, bool)

	//config-unmarshaler
	ConfigUnmarshaler() config.IUnmarshaler

	//message queue
	Nats() *nats.Conn
}
