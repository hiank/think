package think

import (
	"context"

	"github.com/hiank/think/db"
	"github.com/hiank/think/doc/file"
	"github.com/nats-io/nats.go"
)

//utilset provide unique utils
//only provide 'Set' method to get
type utilset interface {
	//TODO base context
	TODO() context.Context

	//Dataset read-write game data
	DBS() db.DBS

	//Fat decoder for decode to given values
	Fat() file.Decoder

	//message queue
	Nats() *nats.Conn
}

//FuncKvDialer convert func to db.KvDialer
type FuncKvDialer func(context.Context, ...db.DialOption) (db.KvDB, error)

func (f FuncKvDialer) Dial(ctx context.Context, opts ...db.DialOption) (db.KvDB, error) {
	return f(ctx, opts...)
}
