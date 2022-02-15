package think

import (
	"context"
	"sync"

	"github.com/hiank/think/db"
	"github.com/hiank/think/doc/file"
	"github.com/hiank/think/run"
	"github.com/nats-io/nats.go"
	"k8s.io/klog/v2"
)

//ErrNoAwake `Awake` has not been executed
const ErrNoAwake = run.Err("think: should do `Awake` before")

//makeOptions make options with given Option
//some field maybe contains default value
func makeOptions(opts ...Option) options {
	dopts := options{
		todo:    context.TODO(),
		mdialer: make(map[db.KeyTag]db.KvDialer),
		mdopts:  make(map[db.KeyTag][]db.DialOption),
	}
	for _, opt := range opts {
		opt.apply(&dopts)
	}
	return dopts
}

type uniqueSet struct {
	cancel   context.CancelFunc
	todo     context.Context
	natsconn *nats.Conn
	decoder  file.Decoder
	dbset    db.DBS
}

func (sm *uniqueSet) TODO() context.Context {
	return sm.todo
}

//Decoder get config decoder
func (sm *uniqueSet) Decoder() file.Decoder {
	return sm.decoder
}

//DBS get database set
func (sm *uniqueSet) DBS() db.DBS {
	return sm.dbset
}

//Nats get nats conn
func (sm *uniqueSet) Nats() *nats.Conn {
	return sm.natsconn
}

var (
	unique *uniqueSet
	once   *sync.Once = new(sync.Once)
)

//Awake create unique object with given options
func Awake(opts ...Option) (done bool) {
	once.Do(func() {
		defer func() {
			if !done {
				unique, once = nil, new(sync.Once)
			}
		}()
		dopts := makeOptions(opts...)
		unique = &uniqueSet{decoder: file.Fat()}
		unique.todo, unique.cancel = context.WithCancel(dopts.todo)
		run.TODO(unique.todo)
		var err error
		if dopts.natsUrl != "" {
			if unique.natsconn, err = nats.Connect(dopts.natsUrl); err != nil {
				klog.Warning("nats connect failed: ", err)
			}
		}
		mdb := make(map[db.KeyTag]db.KvDB)
		for kt, dialer := range dopts.mdialer {
			if mdb[kt], err = dialer.Dial(unique.todo, dopts.mdopts[kt]...); err != nil {
				klog.Warning("k-v database dial failed: ", err)
				delete(mdb, kt)
			}
		}
		unique.dbset, done = db.NewDBS(mdb), true
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

//Destroy destroy the unique
//NOTE: if unique not generate, will panic
func Destroy() {
	defer func() {
		unique, once = nil, new(sync.Once)
	}()
	once.Do(func() {
		panic(ErrNoAwake)
	})
	unique.natsconn.Close()
	unique.cancel() //cancel todo
}
