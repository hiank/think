package run

import (
	"container/list"
	"context"
	"io"
	"time"

	"k8s.io/klog/v2"
)

type taskWorker struct {
	wt Task //working Task
	wc chan Task
	l  *list.List
	C  <-chan time.Time
}

func newTaskWorker(ctx context.Context, timeout time.Duration) (tw *taskWorker) {
	tw = &taskWorker{
		wc: make(chan Task, 24),
		l:  list.New(),
	}
	reset := func() {}
	if timeout > 0 {
		ticker := time.NewTicker(timeout)
		tw.C, reset = ticker.C, func() { ticker.Reset(timeout) }
	}
	go tw.process(ctx, reset)
	return
}

func (tw *taskWorker) push(t Task) {
	if tw.l.Len() == 0 && tw.wt == nil {
		select {
		case tw.wc <- t:
		default:
			tw.wt = t
		}
	} else {
		tw.l.PushBack(t)
	}
}

func (tw *taskWorker) process(ctx context.Context, reset func()) {
	for t := range tw.wc {
		select {
		case <-ctx.Done(): //clear wc's cache
		default:
			reset() //latest recv task
			t.Process()
		}
	}
}

func (tw *taskWorker) work() (sc chan<- Task, t Task) {
	if tw.wt == nil {
		if tw.l.Len() == 0 {
			return
		}
		tw.wt = tw.l.Remove(tw.l.Front()).(Task)
	}
	return tw.wc, tw.wt
}

func (tw *taskWorker) close() {
	close(tw.wc)
}

type tasker struct {
	ctx context.Context
	tc  chan<- Task //Task channel for add a new Task
	io.Closer
}

func NewTasker(ctx context.Context, timeout time.Duration) Tasker {
	tc := make(chan Task, 16)
	t := &tasker{
		tc: tc,
	}
	t.ctx, t.Closer = StartHealthyMonitoring(ctx)
	go t.loop(tc, timeout)
	return t
}

func (t *tasker) loop(tc <-chan Task, timeout time.Duration) {
	defer t.Close()
	worker := newTaskWorker(t.ctx, timeout)
	defer worker.close()

	for {
		sc, wt := worker.work()
		select {
		case <-t.ctx.Done():
			return
		case v := <-tc:
			worker.push(v)
		case sc <- wt:
			worker.wt = nil
		case <-worker.C:
			return //Tasker will closed when non task long time
		}
	}
}

func (t *tasker) Add(tk Task) (err error) {
	select {
	case <-t.ctx.Done():
		err = t.ctx.Err()
	default:
		t.tc <- tk
		//tc will only be closed after ctx cancelled
		//so when tc closed, ctx.Done must be respond
	}
	return err
}

func (*tasker) internalOnly() {}

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
