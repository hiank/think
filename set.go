package think

import (
	"context"
	"io"
	"sync"

	"github.com/hiank/think/auth"
	"github.com/hiank/think/defaults"
	"github.com/hiank/think/doc/sys"
	"github.com/hiank/think/run"
	"github.com/hiank/think/store"
	"github.com/nats-io/nats.go"
	"k8s.io/klog/v2"
)

const (
	ErrInvalidInitialize = run.Err("think: invalid initialize. can only be initialized ont the first call")
)

func defaultOptions() options {
	return options{
		todo:    context.TODO(),
		mdb:     make(map[DBTag]DB),
		natsUrl: defaults.NatsUrl(),
	}
}

type uniqueSet struct {
	todo     context.Context
	tkset    auth.Tokenset
	natsconn *nats.Conn
	mdb      map[DBTag]store.EasyDictionary
	fat      *sys.Fat
	io.Closer
}

func (us *uniqueSet) TODO() context.Context {
	return us.todo
}

func (us *uniqueSet) Tokenset() auth.Tokenset {
	return us.tkset
}

func (us *uniqueSet) Sys() *sys.Fat {
	return us.fat
}

func (us *uniqueSet) DB(tag DBTag) (ed store.EasyDictionary, found bool) {
	ed, found = us.mdb[tag]
	return
}

// Nats get nats conn
func (us *uniqueSet) Nats() *nats.Conn {
	return us.natsconn
}

var (
	unique *uniqueSet
	once   sync.Once
)

// Set utils set
// NOTE: it would panic without call 'Init' method to generate an unique object
func Set(opts ...Option) utilset {
	var done bool
	once.Do(func() {
		dopts := defaultOptions()
		for _, opt := range opts {
			opt.Apply(&dopts)
		}
		todo, closer := run.StartHealthyMonitoring(dopts.todo, destroy)
		// one.TODO(todo)
		unique = &uniqueSet{
			todo:     todo,
			tkset:    auth.NewTokenset(todo),
			mdb:      dialDB(todo, dopts.mdb),
			fat:      sys.NewFat(),
			natsconn: dialNats(dopts.natsUrl),
			Closer:   closer,
		}
		done = true
	})
	if !done && len(opts) > 0 {
		panic(ErrInvalidInitialize)
	}
	return unique
}

func dialNats(url string) (conn *nats.Conn) {
	if url != "" {
		var err error
		if conn, err = nats.Connect(url); err != nil {
			klog.Warning("nats connect failed: ", err)
		}
	}
	return
}

func dialDB(ctx context.Context, mdb map[DBTag]DB) (out map[DBTag]store.EasyDictionary) {
	out = make(map[DBTag]store.EasyDictionary)
	if dbcnt := len(mdb); dbcnt > 0 {
		type result struct {
			err  error
			dict store.EasyDictionary
			tag  DBTag
		}
		pp := make(chan result, dbcnt)
		for _, cfg := range mdb {
			go func(ctx context.Context, cfg DB, pp chan<- result) {
				dict, err := cfg.Dialer.Dial(ctx, cfg.Opts...)
				pp <- result{err, dict, cfg.Tag}
			}(ctx, cfg, pp)
		}
		for rlt := range pp {
			if rlt.err != nil {
				klog.Warning("failed dial to database", rlt.err)
			} else {
				out[rlt.tag] = rlt.dict
			}
			if dbcnt--; dbcnt == 0 {
				close(pp)
			}
		}
	}
	return
}

// destroy destroy the unique
func destroy() {
	if unique.natsconn != nil {
		unique.natsconn.Close()
	}
	for _, dict := range unique.mdb {
		dict.Close()
	}
	unique = nil
}
