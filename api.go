package think

import (
	"context"
	"io"

	"github.com/hiank/think/auth"
	"github.com/hiank/think/doc/sys"
	"github.com/hiank/think/store"
	"github.com/hiank/think/store/db"
	"github.com/nats-io/nats.go"
)

// utilset provide unique utils
// only provide 'Set' method to get
type utilset interface {
	//TODO base context
	//one.TODO()
	TODO() context.Context

	//Tokenset one.TokenSet()
	Tokenset() auth.Tokenset

	//DB get cached database
	DB(tag DBTag) (ed store.EasyDictionary, found bool)

	//Sys config unmarshaler
	Sys() *sys.Fat

	//message queue
	Nats() *nats.Conn

	//Close and clean
	io.Closer
}

type DBTag int

type DB struct {
	Tag    DBTag           //tag in set
	Dialer db.Dialer       //dialer
	Opts   []db.DialOption //dial options
}
