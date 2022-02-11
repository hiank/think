package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/db"
	"k8s.io/klog/v2"
)

type liteDB struct {
	ctx   context.Context
	rcli  *redis.Client
	coder db.BytesCoder
}

func (ld *liteDB) Get(k string, v interface{}) (found bool, err error) {
	str, err := ld.rcli.Get(ld.ctx, k).Result()
	if err == nil {
		err, found = ld.coder.Decode([]byte(str), v), true
	}
	return found, ld.updateErr(err)
}

func (ld *liteDB) Set(k string, v interface{}) error {
	bytes, err := ld.coder.Encode(v)
	if err == nil {
		err = ld.rcli.Set(ld.ctx, k, string(bytes), 0).Err()
	}
	return ld.updateErr(err)
}

func (ld *liteDB) Del(k string, outs ...interface{}) (err error) {
	str, err := ld.rcli.GetDel(ld.ctx, k).Result()
	if err == nil {
		for _, out := range outs {
			if terr := ld.coder.Decode([]byte(str), out); terr != nil {
				klog.Warning(terr)
			}
		}
	}
	return ld.updateErr(err)
}

func (ld *liteDB) updateErr(err error) error {
	switch err {
	case redis.Nil:
		err = db.ErrNotFound
	}
	return err
}

func (ld *liteDB) Close() error {
	return ld.rcli.Close()
}
