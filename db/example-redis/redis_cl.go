package main

import (
	"context"

	"github.com/hiank/think/db"
	"k8s.io/klog/v2"
)

func main() {

	klog.InitFlags(nil)
	// klog.Infoln("why no log")

	ctx := context.Background()
	addr := "redis-master:tcp-redis" //k8s.TryServiceURL(ctx, k8s.TypeKubeIn, "redis-master", "tcp-redis")
	klog.Infoln("redis-master addr : ", addr)
	cli, err := db.NewVerifiedRedisCLI(ctx, &db.RedisConf{
		CheckMillisecond: 500,
		TimeoutSecond:    10,
		Addr:             addr,
		Password:         "oEAlfD10gQ",
		DB:               0,
	})
	if err != nil {
		klog.Infoln(err)
		return
	}
	defer cli.Close()

	status := cli.Set(ctx, "testInt", 1, 0)
	if status.Err() != nil {
		klog.Infoln(status.Err())
		return
	}

	val, err := cli.Get(ctx, "testInt").Result()
	if err != nil {
		klog.Infoln(err)
		return
	}
	klog.Infof("result : %v\n", val)

	slaveCli, err := db.NewVerifiedRedisCLI(ctx, &db.RedisConf{
		Password:         "oEAlfD10gQ",
		DB:               0,
		Addr:             "redis-slave:6379",
		TimeoutSecond:    10,
		CheckMillisecond: 500,
	})
	if err != nil {
		klog.Infoln(err)
		return
	}
	defer slaveCli.Close()

	val, err = slaveCli.Get(ctx, "testInt").Result()
	if err != nil {
		klog.Infoln(err)
		return
	}
	klog.Infoln("result from slave: ", val)
}
