package redis

import (
	"context"

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

func NewKvDB(ctx context.Context, opt *redis.Options) db.KvDB {
	return &liteDB{
		ctx:    ctx,
		Client: redis.NewClient(opt),
	}
}
