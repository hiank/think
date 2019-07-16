package pool

import (
	"sync"
	"context"
)


type runtine struct {

	context.Context
	cancel 	context.CancelFunc
}



type runtineHub struct {

	mtx 	sync.RWMutex
	hub 	map[string]*runtine
}

func newRuntineHub() *runtineHub {

	return &runtineHub{
		hub 	: make(map[string]*runtine),
	}
}

func (rb *runtineHub) get(ctx context.Context, key string) *runtine {

	rb.mtx.Lock()
	defer rb.mtx.Unlock()

	r, ok := rb.hub[key]
	if !ok {
		ctx, cancel := context.WithCancel(ctx)
		r = &runtine{
			Context : ctx,
			cancel 	: cancel,
		}
		rb.hub[key] = r
	}
	return r
}

func (rb *runtineHub) delete(key string) {

	rb.mtx.Lock()
	defer rb.mtx.Unlock()

	r, ok := rb.hub[key]
	if ok {
		r.cancel()
		delete(rb.hub, key)
	}
}