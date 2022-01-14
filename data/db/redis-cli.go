package db

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/kube"
)

func defaultRedisOptions() redisOptions {
	return redisOptions{
		writerOpt: kube.RedisMasterOptions(),
		readerOpt: kube.RedisSlaveOptions(),
	}
}

type redisHub struct {
	ctx    context.Context
	cancel context.CancelFunc
	reader *redis.Client
	writer *redis.Client
}

//NewRedisDB create redis client
//NOTE: without RedisOption, it will connect to k8s's redis server
func NewRedisDB(ctx context.Context, opts ...RedisOption) IClient {
	dopts := defaultRedisOptions()
	for _, opt := range opts {
		opt.apply(&dopts)
	}
	hub := &redisHub{}
	if dopts.writerOpt != nil {
		hub.writer = redis.NewClient(dopts.writerOpt)
		hub.reader = hub.writer
	}
	if dopts.readerOpt != nil {
		hub.reader = redis.NewClient(dopts.readerOpt)
		if hub.writer == nil {
			hub.writer = hub.reader
		}
	}
	go func(ctx context.Context) {
		<-ctx.Done()
		hub.Close()
	}(ctx)
	hub.ctx, hub.cancel = context.WithCancel(ctx)
	return hub
}

func (rh *redisHub) HGet(hashKey, fieldKey string) (IParser, error) {
	cmd := rh.reader.HGet(rh.ctx, hashKey, fieldKey)
	return cmd, cmd.Err()
}

func (rh *redisHub) HSet(hashKey string, values ...interface{}) error {
	return rh.writer.HSet(rh.ctx, hashKey, values...).Err()
}

func (rh *redisHub) Close() error {
	rh.cancel()
	err := rh.writer.Close()
	rh.reader.Close()
	return err
}

//******************create options*******************
type redisOptions struct {
	readerOpt *redis.Options
	writerOpt *redis.Options
}

type RedisOption interface {
	apply(*redisOptions)
}

type funcRedisOption func(*redisOptions)

func (fro funcRedisOption) apply(opts *redisOptions) {
	fro(opts)
}

//WithRedisWriterOption set redis option for writer
//NOTE: if opt is nil, writer option will be nil, client will use reader as writer
func WithRedisWriterOption(opt *redis.Options) RedisOption {
	return funcRedisOption(func(opts *redisOptions) {
		opts.writerOpt = opt
	})
}

//WithRedisReaderOption set redis option for reader
//NOTE: if opt is nil, reader option will be nil, client will use writer as reader
func WithRedisReaderOption(opt *redis.Options) RedisOption {
	return funcRedisOption(func(opts *redisOptions) {
		opts.readerOpt = opt
	})
}
