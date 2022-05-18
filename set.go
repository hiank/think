package think

import (
	"context"
	"io"
	"sync"

	"github.com/hiank/think/doc/sys"
	"github.com/hiank/think/kube"
	"github.com/hiank/think/run"
	"github.com/hiank/think/store"
	"github.com/nats-io/nats.go"
	"k8s.io/klog/v2"
)

//ErrNoAwake `Awake` has not been executed
const ErrNoAwake = run.Err("think: should do `Awake` before")

func defaultOptions() options {
	return options{
		todo:    context.TODO(),
		mdb:     make(map[DBTag]DB),
		natsUrl: kube.NatsUrl(),
	}
}

type uniqueSet struct {
	// cancel   context.CancelFunc
	todo     context.Context
	natsconn *nats.Conn
	mdb      map[DBTag]store.EasyDictionary
	fat      *sys.Fat
	io.Closer
}

func (us *uniqueSet) TODO() context.Context {
	return us.todo
}

func (us *uniqueSet) Sys() *sys.Fat {
	return us.fat
}

func (us *uniqueSet) DB(tag DBTag) (ed store.EasyDictionary, found bool) {
	ed, found = us.mdb[tag]
	return
}

//Nats get nats conn
func (us *uniqueSet) Nats() *nats.Conn {
	return us.natsconn
}

var (
	unique *uniqueSet
	once   *sync.Once = new(sync.Once)
)

//Awake create unique object with given options
func Awake(opts ...Option) (done bool) {
	once.Do(func() {
		dopts := defaultOptions()
		for _, opt := range opts {
			opt.Apply(&dopts)
		}
		healthy := run.NewHealthy()
		todo, cancel := context.WithCancel(dopts.todo)
		unique = &uniqueSet{
			mdb:      dialDB(todo, dopts.mdb),
			todo:     todo,
			fat:      sys.NewFat(),
			natsconn: dialNats(dopts.natsUrl),
			Closer:   run.NewHealthyCloser(healthy, cancel),
		}
		run.TODO(todo)
		go healthy.Monitoring(todo, destroy)
		done = true
	})
	return
}

//Set utils set
//NOTE: it would panic without call 'Init' method to generate an unique object
func Set() utilset {
	once.Do(func() {
		once = new(sync.Once)
		panic(ErrNoAwake)
	})
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

//destroy destroy the unique
func destroy() {
	if unique.natsconn != nil {
		unique.natsconn.Close()
	}
	for _, dict := range unique.mdb {
		dict.Close()
	}
	unique, once = nil, new(sync.Once)
}
