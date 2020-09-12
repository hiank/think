package db

import (
	"context"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/net/k8s"
	"github.com/hiank/think/token"
)

//Rdb
const (
	RdbMaster = "redis-master"
	RdbSlave  = "redis-slave"
)

// tryRedisClient new redis db client
// note: 这个函数可能是个耗时函数
func tryRedisClient(ctx context.Context, rdbName string) *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     k8s.TryServiceURL(ctx, k8s.TypeKubIn, rdbName, "redis"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

var _singleRedisMaster *redis.Client
var _singleRedisMasterOnce sync.Once

// TryRedisMaster redis-master in k8s
// k8s没有需要的redis 服务时，会抛出异常
func TryRedisMaster() *redis.Client {

	_singleRedisMasterOnce.Do(func() {
		_singleRedisMaster = tryRedisClient(token.BackgroundLife().Context, RdbMaster)
	})
	return _singleRedisMaster
}

var _singleRedisSlave *redis.Client
var _singleRedisSlaveOnce sync.Once

// TryRedisSlave redis-slave in k8s
func TryRedisSlave() *redis.Client {

	_singleRedisSlaveOnce.Do(func() {
		_singleRedisSlave = tryRedisClient(token.BackgroundLife().Context, RdbSlave)
	})
	return _singleRedisSlave
}
