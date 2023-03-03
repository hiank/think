package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/run"
	"github.com/hiank/think/store"
	"github.com/hiank/think/store/db"
)

//Dial connect to redis database and return connected client or error
func Dial(ctx context.Context, opts ...db.DialOption) (d store.Dictionary[string], err error) {
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
		ctx, closer := run.StartHealthyMonitoring(ctx, run.CloserToDoneHook(cli))
		d = &liteDB{
			ctx:    ctx,
			rcli:   cli,
			coder:  dopts.Coder,
			Closer: closer,
		}
	}
	return
}
