package box

import (
	"context"
	"io"
	"time"

	"github.com/hiank/think/run"
)

const ContextkeyTokenUid run.Contextkey = "contextkey-uid"

type TokenOption run.Option[*token]

type token struct {
	timeout time.Duration
	context.Context
	io.Closer
}

func WithTokenTimeout(timeout time.Duration) TokenOption {
	return run.FuncOption[*token](func(tk *token) {
		tk.timeout = timeout
	})
}

// func WithTokenUid(uid string) TokenOption {
// 	return run.FuncOption[*token]
// }

func NewToken(ctx context.Context, opts ...TokenOption) Token {
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
	tk.Context = ctx
	tk.Closer = run.NewOnceCloser(func() error {
		cancel()
		return nil
	})
	return tk
}

func (tk *token) Fork(opts ...TokenOption) Token {
	return NewToken(tk.Context, opts...)
}

func (tk *token) internalOnly() {}
