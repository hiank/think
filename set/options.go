package set

import (
	"github.com/hiank/think/data"
	"github.com/hiank/think/data/db"
)

type options struct {
	natsUrl string //NOTE: for Nats
	mstore  map[data.KeyTag]db.KvDB
}

type Option interface {
	apply(*options)
}

type funcOption func(*options)

func (fo funcOption) apply(opts *options) {
	fo(opts)
}

//WithMemKvDB memory database store
//for high performance
func WithMemKvDB(store db.KvDB) Option {
	return funcOption(func(opts *options) {
		opts.mstore[data.KTMem] = store
	})
}

//WithDiskvDB disk database store
//for persistent
func WithDiskvDB(store db.KvDB) Option {
	return funcOption(func(opts *options) {
		opts.mstore[data.KTDisk] = store
	})
}

//WithNatsUrl nats url
//NOTE: if url is "", natsconn will be nil
func WithNatsUrl(url string) Option {
	return funcOption(func(opts *options) {
		opts.natsUrl = url
	})
}
