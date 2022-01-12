package set

import (
	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/set/db"
)

type initOptions struct {
	redisOptions map[db.RedisTag]*redis.Options
	natsUrl      string
}

type InitOption interface {
	apply(*initOptions)
}

type funcInitOption func(*initOptions)

func (fio funcInitOption) apply(opts *initOptions) {
	fio(opts)
}

//WithRedisMasterOption set redis master option
//NOTE: if opt is nil, redis master will be nil
func WithRedisMasterOption(opt *redis.Options) InitOption {
	return funcInitOption(func(io *initOptions) {
		if opt == nil {
			delete(io.redisOptions, db.RedisTagMaster)
		} else {
			io.redisOptions[db.RedisTagMaster] = opt
		}
	})
}

//WithRedisSlaveOption set redis slave option
//NOTE: if opt is nil, redis slave will be nil
func WithRedisSlaveOption(opt *redis.Options) InitOption {
	return funcInitOption(func(io *initOptions) {
		if opt == nil {
			delete(io.redisOptions, db.RedisTagSlave)
		} else {
			io.redisOptions[db.RedisTagSlave] = opt
		}
	})
}

//WithNatsUrl nats url
//NOTE: if url is "", natsconn will be nil
func WithNatsUrl(url string) InitOption {
	return funcInitOption(func(io *initOptions) {
		io.natsUrl = url
	})
}

// func WithoutRedisMaster
