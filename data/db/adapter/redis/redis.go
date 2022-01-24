package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/data/db"
	"github.com/hiank/think/data/db/doc"
)

type liteDB struct {
	ctx context.Context
	*redis.Client
	makeDoc func([]byte) doc.Doc
}

// func (ld *liteDB) Instance() interface{} {
// 	return ld.Client
// }

func (ld *liteDB) Get(k string, v interface{}) (found bool, err error) {
	str, err := ld.Client.Get(ld.ctx, k).Result()
	if err == nil {
		doc := ld.makeDoc([]byte(str))
		err, found = doc.Decode(v), true
	}
	return
}

func (ld *liteDB) Set(k string, v interface{}) (err error) {
	doc := ld.makeDoc(nil)
	if err = doc.Encode(v); err == nil {
		err = ld.Client.Set(ld.ctx, k, doc.Val(), 0).Err()
	}
	return
}

func (ld *liteDB) Delete(k string) error {
	return ld.Client.Del(ld.ctx, k).Err()
}

// func (ld *liteDB) Close() error {
// 	return ld.cli.Close()
// }

func NewKvDB(ctx context.Context, makeDoc func([]byte) doc.Doc, opt *redis.Options) db.KvDB {
	return &liteDB{
		ctx:     ctx,
		Client:  redis.NewClient(opt),
		makeDoc: makeDoc,
	}
}
