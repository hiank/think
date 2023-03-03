package auth

import (
	"context"
	"io"
	"sync"

	"github.com/hiank/think/run"
)

const (
	ErrNonRootoken = run.Err("one: non root Token")
)

type tokenSet struct {
	ctx context.Context //base to TODO
	m   sync.Map
	io.Closer
	internal
}

// Derive copyright Token with given options
func (ts *tokenSet) Derive(key string, topts ...TokenOption) Token {
	v, ok := ts.m.Load(key)
	if !ok {
		nt := &token{}
		nt.Context, nt.Closer = run.StartHealthyMonitoring(context.WithValue(ts.ctx, contextkeyToken, key), func() {
			ts.Kill(key)
		})
		v, _ = ts.m.LoadOrStore(key, nt)
	}
	return v.(Token).Fork(topts...)
}

// Kill the given token
func (ts *tokenSet) Kill(key string) (err error) {
	if v, loaded := ts.m.LoadAndDelete(key); loaded {
		v.(Token).Close()
	} else {
		err = ErrNonRootoken
	}
	return
}

func NewTokenset(ctx context.Context) Tokenset {
	ts := &tokenSet{}
	ts.ctx, ts.Closer = run.StartHealthyMonitoring(ctx, func() {
		*ts = tokenSet{}
	})
	return ts
}
