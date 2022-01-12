package db

import (
	"os"

	"github.com/go-redis/redis/v8"
)

const (
	tagMaster = iota
	tagSlave
)

//-----------------set for redis-------------------//
var (
	DefaultMasterOption = &redis.Options{
		Addr:     "redis-master:tcp-redis",
		Password: os.Getenv("redis-password"),
		DB:       0,
	}
	DefaultSlaveOption = &redis.Options{
		Addr:     "redis-slave:tcp-redis",
		Password: os.Getenv("redis-password"),
		DB:       0,
	}
)

type RedisTag int

const (
	RedisTagMaster RedisTag = tagMaster
	RedisTagSlave  RedisTag = tagSlave
)

//-----------------set for nats-------------------//
var (
	DefaultNatsUrl = "nats:tcp-nats"
)
