package set

import (
	"context"

	"github.com/hiank/think/db"
	"github.com/hiank/think/doc/file"
	"github.com/nats-io/nats.go"
)

type Set interface {
	//TODO base context
	TODO() context.Context

	//Dataset read-write game data
	DBS() db.DBS

	//Fat decoder for decode to given values
	Fat() file.Decoder

	//message queue
	Nats() *nats.Conn
}
