package run

import (
	"container/list"
	"context"
	"io"
	"time"

	"k8s.io/klog/v2"
)

type tasker struct {
	ctx context.Context
	tc  chan<- Task //Task channel for add a new Task
	io.Closer
}

func NewTasker(ctx context.Context, timeout time.Duration) Tasker {
	ctx, cancel := context.WithCancel(ctx)
	healthy, tc := NewHealthy(), make(chan Task, 16)
	t := &tasker{
		ctx:    ctx,
		tc:     tc,
		Closer: NewHealthyCloser(healthy, cancel),
	}
	go healthy.Monitoring(ctx, func() {
		close(t.tc)
		t.tc = nil //avoid `send on closed channel`
	})
	go t.loop(tc, timeout)
	return t
}

func (t *tasker) loop(tc <-chan Task, timeout time.Duration) {
	wc := make(chan Task, 24) //work chan
	defer close(wc)
	go func(wc <-chan Task) {
		for task := range wc {
			task.Process()
		}
	}(wc)

	var wt Task
	var sc chan<- Task
	ticker, l := time.NewTicker(timeout), list.New()
	for {
		select {
		case v, ok := <-tc:
			if !ok {
				return //tc closed only after ctx cancelled
			}
			if wt == nil {
				wt, sc = v, wc
			} else {
				l.PushBack(v)
			}
		case sc <- wt:
			if em := l.Front(); em != nil {
				wt = l.Remove(em).(Task)
			} else {
				wt, sc = nil, nil
			}
		case <-ticker.C:
			t.Close()
			return //Tasker will closed when non task long time
		}
		ticker.Reset(timeout)
	}
}

func (t *tasker) Add(tk Task) (err error) {
	select {
	case <-t.ctx.Done():
		err = t.ctx.Err()
	case t.tc <- tk:
		//tc will only be closed after ctx cancelled
		//so when tc closed, ctx.Done must be respond
	}
	return err
}

type liteTask[T any] struct {
	v     T
	h     func(T) error
	pperr chan<- error
}

//Process
func (tt *liteTask[T]) Process() (err error) {
	if err = tt.h(tt.v); err != nil && tt.pperr != nil {
		select {
		case tt.pperr <- err:
		case <-time.NewTicker(time.Second).C:
			klog.Warning("cannot send error for invalid 'handle' in long try")
		}
	}
	return
}

func NewLiteTask[T any](h func(T) error, v T, pperr ...chan<- error) Task {
	lt := &liteTask[T]{
		h: h,
		v: v,
	}
	if len(pperr) > 0 {
		lt.pperr = pperr[0]
	}
	return lt
}
