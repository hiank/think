package set

import (
	"errors"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/config"
)

func defaultInitOptions() initOptions {
	return initOptions{
		redisMasterOption: &redis.Options{
			Addr:     "redis-master:tcp-redis",
			Password: os.Getenv("redis-password"),
			DB:       0,
		},
		redisSlaveOption: &redis.Options{
			Addr:     "redis-slave:tcp-redis",
			Password: os.Getenv("redis-password"),
			DB:       0,
		},
	}
}

const (
	redisMasterKey = iota
	redisSlaveKey
)

type getter struct {
	rdbm map[int]*redis.Client
}

func (sm *getter) RedisMasterCli() (cli *redis.Client, ok bool) {
	cli, ok = sm.rdbm[redisMasterKey]
	return
}

func (sm *getter) RedisSlaveCli() (cli *redis.Client, ok bool) {
	cli, ok = sm.rdbm[redisSlaveKey]
	return
}

func (sm *getter) ConfigUnmarshaler() config.IUnmarshaler {
	return config.NewUnmarshaler()
}

var (
	instance *getter
	once     sync.Once
)

//Instance IOpenApi singleton
//NOTE: the parms only take effect in the first call the method. they are used to init instance
func Instance(opts ...InitOption) IOpenApi {
	once.Do(func() {
		dopts := defaultInitOptions()
		for _, opt := range opts {
			opt.apply(&dopts)
		}
		instance = &getter{
			rdbm: make(map[int]*redis.Client),
		}
		if dopts.redisMasterOption != nil {
			instance.rdbm[redisMasterKey] = redis.NewClient(dopts.redisMasterOption)
		}
		if dopts.redisSlaveOption != nil {
			instance.rdbm[redisSlaveKey] = redis.NewClient(dopts.redisSlaveOption)
		}
	})
	return instance
}

//Relase relase the instance
//NOTE: if instance not generate, will panic
func Release() {
	once.Do(func() {
		panic(errors.New("instance not generate now. should not call Release"))
	})
	for _, cli := range instance.rdbm {
		cli.Close()
	}
	instance = nil
	once = sync.Once{}
}
