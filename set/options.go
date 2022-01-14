package set

import (
	"github.com/hiank/think/data/db"
)

type options struct {
	// redisOptions map[db.RedisTag]*redis.Options
	natsUrl  string     //NOTE: for Nats
	memoryDB db.IClient //NOTE: for Dataset
	diskDB   db.IClient //NOTE: for Dataset
}

type Option interface {
	apply(*options)
}

type funcOption func(*options)

func (fo funcOption) apply(opts *options) {
	fo(opts)
}

// //WithRedisMasterOption set redis master option
// //NOTE: if opt is nil, redis master will be nil
// func WithRedisMasterOption(opt *redis.Options) InitOption {
// 	return funcInitOption(func(io *initOptions) {
// 		if opt == nil {
// 			delete(io.redisOptions, db.RedisTagMaster)
// 		} else {
// 			io.redisOptions[db.RedisTagMaster] = opt
// 		}
// 	})
// }

// //WithRedisSlaveOption set redis slave option
// //NOTE: if opt is nil, redis slave will be nil
// func WithRedisSlaveOption(opt *redis.Options) InitOption {
// 	return funcInitOption(func(io *initOptions) {
// 		if opt == nil {
// 			delete(io.redisOptions, db.RedisTagSlave)
// 		} else {
// 			io.redisOptions[db.RedisTagSlave] = opt
// 		}
// 	})
// }

//WithDatasetMemoryDB memory database for Dataset
//NOTE: memory database must set
func WithDatasetMemoryDB(cli db.IClient) Option {
	return funcOption(func(opts *options) {
		opts.memoryDB = cli
	})
}

//WithDatasetDiskDB disk database for Dataset
func WithDatasetDiskDB(cli db.IClient) Option {
	return funcOption(func(opts *options) {
		opts.diskDB = cli
	})
}

//WithNatsUrl nats url
//NOTE: if url is "", natsconn will be nil
func WithNatsUrl(url string) Option {
	return funcOption(func(opts *options) {
		opts.natsUrl = url
	})
}

// func WithoutRedisMaster
