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
	mux.RLock()
	defer mux.RUnlock()
	if len(ctx) > 0 {
		todo = ctx[0]
	}
	return todo
}
