package set

import (
	"errors"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/config"
	"github.com/hiank/think/set/db"
	"github.com/nats-io/nats.go"
)

func defaultInitOptions() initOptions {
	return initOptions{
		redisOptions: map[db.RedisTag]*redis.Options{db.RedisTagMaster: db.DefaultMasterOption, db.RedisTagSlave: db.DefaultSlaveOption},
		natsUrl:      db.DefaultNatsUrl,
	}
}

type getter struct {
	rdbm     map[db.RedisTag]*redis.Client
	natsconn *nats.Conn
	cfgum    config.IParser
}

//ConfigParser get config parser
func (sm *getter) ConfigParser() config.IParser {
	return sm.cfgum
}

//RedisCli get redis client by given RedisTag (RedisTagMaster || RedisTagSlave)
func (sm *getter) RedisCli(tag db.RedisTag) (cli *redis.Client, ok bool) {
	cli, ok = sm.rdbm[tag]
	return
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
func Init(opts ...InitOption) (done bool) {
	once.Do(func() {
		dopts := defaultInitOptions()
		for _, opt := range opts {
			opt.apply(&dopts)
		}
		unique = &getter{
			rdbm:  make(map[db.RedisTag]*redis.Client),
			cfgum: config.NewParser(),
		}
		for tag, opt := range dopts.redisOptions {
			unique.rdbm[tag] = redis.NewClient(opt)
		}
		if dopts.natsUrl != "" {
			unique.natsconn, _ = nats.Connect(dopts.natsUrl)
		}
		done = true
	})
	return
}

//Unique IOpenApi singleton
//NOTE: it would panic without call 'Init' method to generate an unique object
func Unique() IOpenApi {
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
	for _, cli := range unique.rdbm {
		cli.Close()
	}
	unique.natsconn.Close()

	unique = nil
	once = sync.Once{}
}
