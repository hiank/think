package think

import (
	"context"

	"github.com/hiank/think/run"
	"k8s.io/klog/v2"
)

type options struct {
	natsUrl string //NOTE: for Nats
	todo    context.Context
	mdb     map[DBTag]DB
}

type Option run.Option[*options]

//WithNatsUrl nats url
//NOTE: if url is "", natsconn will be nil
func WithNatsUrl(url string) Option {
	return run.FuncOption[*options](func(opts *options) {
		opts.natsUrl = url
	})
}

//WithTODO base Context
//the todo will cancel when Destroy
func WithTODO(ctx context.Context) Option {
	return run.FuncOption[*options](func(opts *options) {
		opts.todo = ctx
	})
}

//WithDB DB use for dial to database and cache in set
//cfg.Tag use for cache in set
func WithDB(cfg DB) Option {
	return run.FuncOption[*options](func(opts *options) {
		if _, ok := opts.mdb[cfg.Tag]; ok {
			klog.Warningf("tag %d for DB is existed", cfg.Tag)
		}
		opts.mdb[cfg.Tag] = cfg
	})
}
