package run

import (
	"context"
	"sync"
	"time"
)

type ForkOptions struct {
	Timeout time.Duration
}

type CtxKey string

const (
	CtxKeyUid CtxKey = "uid"
)

// type

type tokenSet struct {
	ctx context.Context //base to TODO
	m   sync.Map
	// mux sync.RWMutex
}

func (ts *tokenSet) Get(uid uint64) Token {
	if v, ok := ts.m.Load(uid); ok {
		return v.(Token)
	}
	ctx, tk := context.WithValue(ts.ctx, CtxKeyUid, uid), &token{}
	tk.Context, tk.cancel = context.WithCancel(ctx)

	v, _ := ts.m.LoadOrStore(uid, tk)
	return v.(Token)
}

type token struct {
	context.Context
	cancel context.CancelFunc
}

func (tk *token) Fork(opt ForkOptions) Token {
	ftk := &token{}
	if opt.Timeout > 0 {
		ftk.Context, ftk.cancel = context.WithTimeout(tk.Context, opt.Timeout)
	} else {
		ftk.Context, ftk.cancel = context.WithCancel(tk.Context)
	}
	return ftk
}

func (tk *token) Cancel() {
	tk.cancel()
}
