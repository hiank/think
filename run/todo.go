package run

import (
	"context"
	"sync"
)

var (
	todo = context.TODO()
	mux  sync.RWMutex
)

func TODO(ctx ...context.Context) context.Context {
	if len(ctx) > 0 {
		mux.Lock()
		defer mux.Unlock()
		todo = ctx[0]
		return todo
	}
	mux.RLock()
	defer mux.RUnlock()
	return todo
}
