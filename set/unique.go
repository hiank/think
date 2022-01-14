package set

import (
	"context"
	"errors"
	"sync"

	dset "github.com/hiank/think/data"
	"github.com/hiank/think/fp"
	"github.com/hiank/think/kube"
	"github.com/nats-io/nats.go"
)

func defaultOptions() options {
	return options{
		// redisOptions: map[db.RedisTag]*redis.Options{db.RedisTagMaster: db.DefaultMasterOption, db.RedisTagSlave: db.DefaultSlaveOption},
		natsUrl: kube.NatsUrl(),
	}
}

type getter struct {
	Cancel   context.CancelFunc
	ctx      context.Context
	natsconn *nats.Conn
	textp    fp.IParser
	dataset  dset.IDataset
}

//TextParser get config parser
func (sm *getter) TextParser() fp.IParser {
	return sm.textp
}

// //RedisCli get redis client by given RedisTag (RedisTagMaster || RedisTagSlave)
// func (sm *getter) RedisCli(tag db.RedisTag) (cli *redis.Client, ok bool) {
// 	cli, ok = sm.rdbm[tag]
// 	return
// }

func (sm *getter) Dataset() dset.IDataset {
	return sm.dataset
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
			dataset: dset.NewDataset(ctx, dset.WithMemoryDB(dopts.memoryDB), dset.WithDiskDB(dopts.diskDB)),
			textp:   fp.NewParser(),
		}
		if dopts.natsUrl != "" {
			unique.natsconn, _ = nats.Connect(dopts.natsUrl)
		}
		done = true
	})
	return
}

//Unique ISet singleton
//NOTE: it would panic without call 'Init' method to generate an unique object
func Unique() ISet {
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
