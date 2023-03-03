package redis

import (
	"context"
	"io"

	"github.com/go-redis/redis/v8"
	"github.com/hiank/think/doc"
	"github.com/hiank/think/run"
	"k8s.io/klog/v2"
)

const (
	ErrNotFound = run.Err("redis: not found")
)

type liteDB struct {
	ctx   context.Context
	rcli  *redis.Client
	coder doc.Coder
	io.Closer
}

func (ld *liteDB) Scan(k string, v any) (found bool, err error) {
	str, err := ld.rcli.Get(ld.ctx, k).Result()
	if err == nil {
		ld.coder.Encode([]byte(str))
		err, found = ld.coder.Decode(v), true
	}
	return found, ld.updateErr(err)
}

func (ld *liteDB) Set(k string, v any) (err error) {
	if err = ld.coder.Encode(v); err == nil {
		err = ld.rcli.Set(ld.ctx, k, string(ld.coder.Bytes()), 0).Err()
	}
	return ld.updateErr(err)
}

func (ld *liteDB) Del(k string, outs ...any) (err error) {
	str, err := ld.rcli.GetDel(ld.ctx, k).Result()
	if err == nil {
		if err = ld.coder.Encode([]byte(str)); err == nil {
			for _, out := range outs {
				if terr := ld.coder.Decode(out); terr != nil {
					klog.Warning(terr)
				}
			}
		}
	}
	return ld.updateErr(err)
}

func (ld *liteDB) updateErr(err error) error {
	switch err {
	case redis.Nil:
		err = ErrNotFound
	}
	return err
}
