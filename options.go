package think

import (
	"context"

	"github.com/hiank/think/db"
)

type options struct {
	natsUrl string //NOTE: for Nats
	mdialer map[db.KeyTag]db.KvDialer
	mdopts  map[db.KeyTag][]db.DialOption
	todo    context.Context
}

type Option interface {
	apply(*options)
}

type funcOption func(*options)

func (fo funcOption) apply(opts *options) {
	fo(opts)
}

//WithMemKvDialer memory k-v database dialer and dial options
//for use todo context, delay call until after Awake
func WithMemKvDialer(dialer db.KvDialer, dopts ...db.DialOption) Option {
	return funcOption(func(opts *options) {
		opts.mdialer[db.KTMem] = dialer
		opts.mdopts[db.KTMem] = dopts
	})
}

//WithDiskvDialer disk k-v database dialer and dial options
//for use todo context, delay call until after Awake
func WithDiskvDialer(dialer db.KvDialer, dopts ...db.DialOption) Option {
	return funcOption(func(opts *options) {
		opts.mdialer[db.KTDisk] = dialer
		opts.mdopts[db.KTDisk] = dopts
	})
}

//WithNatsUrl nats url
//NOTE: if url is "", natsconn will be nil
func WithNatsUrl(url string) Option {
	return funcOption(func(opts *options) {
		opts.natsUrl = url
	})
}

//WithTODO base Context
//the todo will cancel when Destroy
func WithTODO(ctx context.Context) Option {
	return funcOption(func(opts *options) {
		opts.todo = ctx
	})
}
