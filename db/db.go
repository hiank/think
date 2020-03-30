package db

import (
	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"github.com/hiank/think/net/k8s"
)

func DialToRedis() {

	addr, err := k8s.ServiceNameWithPort(k8s.TypeKubIn, "redis-master", "redis")
	if err != nil {

		glog.Error(err)
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr 		: addr,
		Password	: "",
		DB 			: 0,
	})
	pong, err := client.Ping().Result()
	glog.Infoln(pong, err)
}