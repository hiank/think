package set

import (
	"context"

	"github.com/hiank/think/data"
	"github.com/hiank/think/fp"
	"github.com/nats-io/nats.go"
)

type ISet interface {
	//TODO base context
	TODO() context.Context

	//Dataset read-write game data
	Dataset() data.IDataset

	//text-parser
	TextParser() fp.IParser

	//message queue
	Nats() *nats.Conn
}
