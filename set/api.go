package set

import (
	dset "github.com/hiank/think/data"
	"github.com/hiank/think/fp"
	"github.com/nats-io/nats.go"
)

type ISet interface {
	// //db-redis
	// RedisCli(db.RedisTag) (*redis.Client, bool)

	//Dataset read-write game data
	Dataset() dset.IDataset

	//text-parser
	TextParser() fp.IParser

	//message queue
	Nats() *nats.Conn
}
