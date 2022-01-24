package set

import (
	"context"

	"github.com/hiank/think/data"
	"github.com/hiank/think/fp"
	"github.com/nats-io/nats.go"
)

type Set interface {
	//TODO base context
	TODO() context.Context

	//Dataset read-write game data
	Dataset() data.Dataset

	//text-parser
	TextParser() fp.Parser

	//message queue
	Nats() *nats.Conn
}
