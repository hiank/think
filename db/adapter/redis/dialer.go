package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/db"
)

var (
	Dialer db.KvDialer = dialer(0)
)

type dialer byte

//Dial connect to redis database and return connected client or error
func (d dialer) Dial(ctx context.Context, opts ...db.DialOption) (kv db.KvDB, err error) {
	dopts := db.DialOptions(opts...)
	var dbVal int64
	if dbVal, err = strconv.ParseInt(dopts.DB, 10, 32); err != nil {
		return nil, fmt.Errorf("DB for redis.Options should be int value: %s", dopts.DB)
	}
	cli := redis.NewClient(&redis.Options{
		DB:          int(dbVal),
		Username:    dopts.Account,
		Password:    dopts.Password,
		DialTimeout: dopts.DialTimeout,
		Addr:        dopts.Addr,
	})
	if err = cli.Ping(ctx).Err(); err == nil {
		kv = &liteDB{
			ctx:  ctx,
			rcli: cli,
		}
	}
	return
}
