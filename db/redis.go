package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

//RedisConf redis config
type RedisConf struct {
	CheckMillisecond int    `json:"redis.CheckMillisecond"` //NOTE: redis 检查间隔，开始时，如果ping redis出错，会间隔此时间再检查一次，单位为毫秒
	TimeoutSecond    int    `json:"redis.TimeoutSecond"`    //NOTE: redis 开始阶段，超过此时长未能连入，则连接失败
	Addr             string //`json:"redis.Addr"`             //NOTE: redis url
	Password         string `json:"redis.Password"` //NOTE: redis 密码
	DB               int    `json:"redis.DB"`       //NOTE: redis 参数
}

//NewVerifiedRedisCLI 获取一个验证过的redis.Client，如果无法连接，返回错误
func NewVerifiedRedisCLI(ctx context.Context, conf *RedisConf) (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(conf.TimeoutSecond))
	defer cancel()
	interval := time.Millisecond * time.Duration(conf.CheckMillisecond)
	for {
		err := cli.Ping(ctx).Err()
		if err == nil {
			return cli, nil
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context done: %v : %v", ctx.Err(), err)
		case <-time.After(interval): //NOTE: 每隔一定时间尝试连接redis
		}
	}
}
