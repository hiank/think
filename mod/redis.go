package mod

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think"
	"github.com/hiank/think/db"
	"github.com/hiank/think/mod/modex"
	"github.com/hiank/think/set"
)

var (
	RedisCLIMaster = &redisCLI{addrGetter: new(RedisAddrMaster)}
	RedisCLISlave  = &redisCLI{addrGetter: new(RedisAddrSlave)}
)

var defaultRedisConf = `{
	"redis.CheckMillisecond": 500,
	"redis.TimeoutSecond": 10,
	"redis.Password": "env:REDIS_PASSWORD"
}`

var defaultRedisAddrConf = `{
	"redis.AddrMaster": {
		"Addr": "redis-master:tcp-redis"
	},
	"redis.AddrSlave": {
		"Addr": "redis-slave:tcp-redis"
	}
}`

type RedisAddrMaster struct {
	*modex.Addr `json:"redis.AddrMaster"`
}

func (master *RedisAddrMaster) Get() *modex.Addr {
	return master.Addr
}

type RedisAddrSlave struct {
	*modex.Addr `json:"redis.AddrSlave"`
}

func (slave *RedisAddrSlave) Get() *modex.Addr {
	return slave.Addr
}

type redisCLI struct {
	*redis.Client
	conf       *db.RedisConf
	addrGetter modex.AddrGetter

	think.IgnoreOnDestroy
}

func (cli *redisCLI) Depend() []think.Module {
	return []think.Module{Config, KubesetIn}
}

//OnCreate 此阶段，需要把配置数据注册到ConfigMod
//NOTE: 使用前，务必完成单元测试运行，确保'default' value合法
func (cli *redisCLI) OnCreate(ctx context.Context) (err error) {
	cli.conf = new(db.RedisConf)
	json.Unmarshal([]byte(defaultRedisConf), cli.conf)
	json.Unmarshal([]byte(defaultRedisAddrConf), cli.addrGetter)
	Config.SignUpValue(set.JSON, cli.conf)
	Config.SignUpValue(set.JSON, cli.addrGetter)
	return
}

func (cli *redisCLI) OnStart(ctx context.Context) (err error) {
	if cli.conf.Addr, err = modex.ParseAddr(cli.addrGetter.Get().Value, KubesetIn); err == nil {
		prefix := "env:"
		if strings.Index(cli.conf.Password, prefix) == 0 {
			cli.conf.Password = os.Getenv(cli.conf.Password[len(prefix):])
		}
		cli.Client, err = db.NewVerifiedRedisCLI(ctx, cli.conf)
	}
	return
}

func (cli *redisCLI) OnStop() {
	cli.Client.Close()
}
