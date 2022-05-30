package one

import (
	"context"
	"sync"

	"github.com/hiank/think/run"
)

var (
	todo     = context.TODO()
	todoOnce sync.Once

	ErrInvalidtodoSet = run.Err("one: invalid todo context set (can only set at the first call)")
)

//TODO root context in frame
func TODO(ctxs ...context.Context) context.Context {
	var done bool
	todoOnce.Do(func() {
		if len(ctxs) > 0 {
			todo = ctxs[0]
		}
		done = true
	})
	if !done && len(ctxs) > 0 {
		panic(ErrInvalidtodoSet)
	}
	return todo
}
