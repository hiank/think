package set

import "github.com/go-redis/redis/v8"

type initOptions struct {
	redisMasterOption *redis.Options
	redisSlaveOption  *redis.Options
	natsUrl           string
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
		io.redisMasterOption = opt
	})
}

//WithRedisSlaveOption set redis slave option
//NOTE: if opt is nil, redis slave will be nil
func WithRedisSlaveOption(opt *redis.Options) InitOption {
	return funcInitOption(func(io *initOptions) {
		io.redisSlaveOption = opt
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
