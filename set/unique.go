package set

import (
	"context"
	"errors"
	"sync"

	"github.com/hiank/think/db"
	"github.com/hiank/think/doc/file"
	"github.com/hiank/think/kube"
	"github.com/nats-io/nats.go"
)

func defaultOptions() options {
	return options{
		natsUrl: kube.NatsUrl(),
		mstore:  make(map[db.KeyTag]db.KvDB),
	}
}

type getter struct {
	Cancel   context.CancelFunc
	ctx      context.Context
	natsconn *nats.Conn
	decoder  file.Decoder
	dbset    db.DBS
}

func (sm *getter) TODO() context.Context {
	return sm.ctx
}

//Fat get config decoder
func (sm *getter) Fat() file.Decoder {
	return sm.decoder
}

func (sm *getter) DBS() db.DBS {
	return sm.dbset
}

//Nats get nats conn
func (sm *getter) Nats() *nats.Conn {
	return sm.natsconn
}

var (
	unique *getter
	once   sync.Once
)

//Init create unique object with given options
func Init(opts ...Option) (done bool) {
	once.Do(func() {
		dopts := defaultOptions()
		for _, opt := range opts {
			opt.apply(&dopts)
		}
		ctx, cancel := context.WithCancel(context.Background())
		unique = &getter{
			ctx:     ctx,
			Cancel:  cancel,
			dbset:   db.NewDBS(dopts.mstore),
			decoder: file.Fat(),
		}
		if dopts.natsUrl != "" {
			unique.natsconn, _ = nats.Connect(dopts.natsUrl)
		}
		done = true
	})
	return
}

//Unique Set singleton
//NOTE: it would panic without call 'Init' method to generate an unique object
func Unique() Set {
	once.Do(func() {
		panic(errors.New("unique not generate now. you should call 'set.Init' to generate an unique object"))
	})
	return unique
}

//Clear clear the unique
//NOTE: if unique not generate, will panic
func Clear() {
	once.Do(func() {
		panic(errors.New("unique not generate now. should not call Release"))
	})

	unique.natsconn.Close()
	unique.Cancel() //cancel ctx

	unique = nil
	once = sync.Once{}
}
