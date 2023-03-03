package auth

import (
	"context"
	"io"
	"time"

	"github.com/hiank/think/run"
)

const (
	contextkeyToken run.Contextkey = "contextkey-token"
)

type TokenOption run.Option[*token]

type token struct {
	timeout time.Duration
	context.Context
	io.Closer
	internal
}

func WithTokenTimeout(timeout time.Duration) TokenOption {
	return run.FuncOption[*token](func(tk *token) {
		tk.timeout = timeout
	})
}

func newToken(ctx context.Context, opts ...TokenOption) Token {
	tk := &token{}
	for _, opt := range opts {
		opt.Apply(tk)
	}
	var cancel context.CancelFunc
	if tk.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, tk.timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	tk.Context, tk.Closer = ctx, run.NewOnceCloser(func() error {
		cancel()
		return nil
	})
	return tk
}

func (tk *token) Fork(opts ...TokenOption) Token {
	return newToken(tk.Context, opts...)
}

func (tk *token) ToString() (key string) {
	if v, ok := tk.Value(contextkeyToken).(string); ok {
		key = v
	}
	return
}

// func (tk *token) internalOnly() {}
