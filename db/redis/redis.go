package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/db"
)

type liteDB struct {
	ctx context.Context
	*redis.Client
	coder db.BytesCoder
}

func (ld *liteDB) Get(k string, v interface{}) (found bool, err error) {
	str, err := ld.Client.Get(ld.ctx, k).Result()
	if err == nil {
		err, found = ld.coder.Decode([]byte(str), v), true
	}
	return
}

func (ld *liteDB) Set(k string, v interface{}) (err error) {
	bytes, err := ld.coder.Encode(v)
	if err == nil {
		err = ld.Client.Set(ld.ctx, k, string(bytes), 0).Err()
	}
	return
}

func (ld *liteDB) Delete(k string) error {
	return ld.Client.Del(ld.ctx, k).Err()
}

//Dial connect to redis database and return connected client or error
func Dial(ctx context.Context, opts ...db.DialOption) (kv db.KvDB, err error) {
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
			ctx:    ctx,
			Client: cli,
		}
	}
	return
}
