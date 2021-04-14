package mod

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think"
	"github.com/hiank/think/db"
	"github.com/hiank/think/set"
)

var (
	RedisCLIMaster = &redisCLI{conf: new(db.RedisConf), defaultAddr: "redis-master"}
	RedisCLISlave  = &redisCLI{conf: new(db.RedisConf), defaultAddr: "redis-slave"}
)

var defaultRedisConf = `{
	"redis.CheckMillisecond": 500,
	"redis.TimeoutSecond": 10,
	"redis.Addr": %s,
	"redis.Password": "env:REDIS_PASSWORD"
}`

type redisCLI struct {
	*redis.Client
	conf        *db.RedisConf
	defaultAddr string //NOTE: 默认addr

	think.IgnoreOnStop
	think.IgnoreOnDestroy
}

func (cli *redisCLI) Depend() []think.Module {
	return []think.Module{Config}
}

//OnCreate 此阶段，需要把配置数据注册到ConfigMod
func (cli *redisCLI) OnCreate(ctx context.Context) error {
	strConf := fmt.Sprintf(defaultRedisConf, cli.defaultAddr)
	json.Unmarshal([]byte(strConf), cli.conf)
	Config.SignUpValue(set.JSON, cli.conf)
	return nil
}

func (cli *redisCLI) OnStart(ctx context.Context) (err error) {
	prefix := "env:"
	if strings.Index(cli.conf.Password, prefix) == 0 {
		cli.conf.Password = os.Getenv(cli.conf.Password[len(prefix):])
	}
	cli.Client, err = db.NewVerifiedRedisCLI(ctx, cli.conf)
	return err
}
