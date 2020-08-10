package db

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
	"github.com/hiank/think/net/k8s"
	"github.com/hiank/think/token"
)

//Rdb
const (
	RdbMaster = "redis-master"
	RdbSlave  = "redis-slave"
)

// RedisClient redis db client
type RedisClient struct {
	ctx context.Context
	rdb *redis.Client
}

// newRedisClient new redis db client
// note: 这个函数可能是个耗时函数
func newRedisClient(ctx context.Context, rdbName string) *RedisClient {

	addr, err := k8s.ServiceNameWithPort(context.Background(), k8s.TypeKubIn, rdbName, "redis")
	if err != nil {
		glog.Error(err)
		return nil
	}

	return &RedisClient{
		ctx: ctx,
		rdb: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		}),
	}
}

//Get get value by key
func (rc *RedisClient) Get(key string) *redis.StringCmd {

	return rc.rdb.Get(rc.ctx, key)
}

//Set set value with key
func (rc *RedisClient) Set(key string, val interface{}) *redis.StatusCmd {

	return rc.rdb.Set(rc.ctx, key, val, 0)
}

var _singleRedisMaster *RedisClient
var _singleRedisMasterOnce sync.Once

// RedisMaster redis-master in k8s
func RedisMaster() *RedisClient {

	_singleRedisMasterOnce.Do(func() {
		_singleRedisMaster = newRedisClient(token.BackgroundLife().Context, RdbMaster)
		go func() {
			<-_singleRedisMaster.ctx.Done()
			_singleRedisMaster = nil
			_singleRedisMasterOnce = sync.Once{}
		}()
	})
	return _singleRedisMaster
}

var _singleRedisSlave *RedisClient
var _singleRedisSlaveOnce sync.Once

// RedisSlave redis-slave in k8s
func RedisSlave() *RedisClient {

	_singleRedisSlaveOnce.Do(func() {
		_singleRedisSlave = newRedisClient(token.BackgroundLife().Context, RdbSlave)
		go func() {
			<-_singleRedisSlave.ctx.Done()
			_singleRedisSlave = nil
			_singleRedisSlaveOnce = sync.Once{}
		}()
	})
	return _singleRedisSlave
}
