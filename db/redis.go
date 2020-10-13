package db

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

//RedisConf redis config
type RedisConf struct {
	SlaveURL  string `json:"redis.SlaveURL"`  //NOTE: slave url
	MasterURL string `json:"redis.MasterURL"` //NOTE: master url
	Password  string `json:"redis.Password"`  //NOTE: redis 密码
	DB        int    `json:"redis.DB"`        //NOTE: redis 参数
}

//AutoRedis 维护redis client
type AutoRedis struct {
	context.Context
	rc         *RedisConf
	slave      *redis.Client
	slaveOnce  *sync.Once
	master     *redis.Client
	masterOnce *sync.Once
}

//NewAutoRedis 创建新的AutoRedis
func NewAutoRedis(ctx context.Context, rc *RedisConf) *AutoRedis {
	return &AutoRedis{
		Context:    ctx,
		rc:         rc,
		slaveOnce:  new(sync.Once),
		masterOnce: new(sync.Once),
	}
}

//TryMaster try to get client connectted redis-master
func (ar *AutoRedis) TryMaster() *redis.Client {

	ar.masterOnce.Do(func() {
		ar.master = ar.tryClient(&redis.Options{
			Addr:     ar.rc.MasterURL,
			Password: ar.rc.Password,
			DB:       ar.rc.DB,
		}, &ar.masterOnce)
	})
	return ar.master
}

//TrySlave try to get client connectted redis-slave
func (ar *AutoRedis) TrySlave() *redis.Client {
	ar.slaveOnce.Do(func() {
		ar.slave = ar.tryClient(&redis.Options{
			Addr:     ar.rc.SlaveURL,
			Password: ar.rc.Password,
			DB:       ar.rc.DB,
		}, &ar.slaveOnce)
	})
	return ar.slave
}

func (ar *AutoRedis) connected(rc *redis.Client) <-chan error {

	failed := make(chan error)
	go func() {

		timeout, interval := time.Second*30, time.Millisecond*300
		for {
			err := rc.Ping(ar.Context).Err()
			if err == nil {
				close(failed)
				return
			}
			select {
			case <-ar.Context.Done():
				failed <- errors.New("Context Done")
				return
			case <-time.After(timeout):
				failed <- errors.New("timeout")
				return
			case <-time.After(interval): //NOTE: 每隔一定时间尝试连接redis
			}
		}
	}()
	return failed
}

func (ar *AutoRedis) tryClient(opt *redis.Options, once **sync.Once) *redis.Client {

	client := redis.NewClient(opt)
	if err := <-ar.connected(client); err != nil {
		*once = new(sync.Once)
		panic(err)
	}
	return client
}

func (ar *AutoRedis) syncNewConnectedClient(opt *redis.Options) (*redis.Client, error) {

	client := redis.NewClient(opt)
	if err := <-ar.connected(client); err != nil {
		return nil, err
	}
	return client, nil
}
