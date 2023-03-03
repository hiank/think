package run

import (
	"container/list"
	"context"
	"io"
	"time"

	"k8s.io/klog/v2"
)

const (
	ErrUnrecoverable = Err("unrecoverable error")
)

type taskWorker struct {
	wt Task //working Task
	wc chan Task
	l  *list.List
	C  <-chan time.Time
	D  <-chan bool //unrecoverable error chan
}

func newTaskWorker(ctx context.Context, timeout time.Duration) (tw *taskWorker) {
	tw = &taskWorker{
		wc: make(chan Task, 24),
		l:  list.New(),
	}
	reset, dis := func() {}, make(chan bool)
	if timeout > 0 {
		ticker := time.NewTicker(timeout)
		tw.C, reset = ticker.C, func() { ticker.Reset(timeout) }
	}
	go tw.process(ctx, reset, dis)
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

func (tw *taskWorker) process(ctx context.Context, reset func(), dis chan<- bool) {
	for t := range tw.wc {
		select {
		case <-ctx.Done(): //clear wc's cache
		default:
			reset() //latest recv task
			if t.Process() == ErrUnrecoverable {
				close(dis)   //close chan after receive an unrecoverable error
				<-ctx.Done() //must trigger the shutdown of the user
			}
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
		case <-worker.D:
			return //worker encoutered an unrecoverable error
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
	v      T
	h      func(T) error
	hooker Hooker[error]
}

//Hook for Hooker[error]
func (*liteTask[T]) Hook(err error) {
	//do nothing for error recved
	klog.Warning("failed to process task:", err)
}

//Process
func (tt *liteTask[T]) Process() (err error) {
	if err = tt.h(tt.v); err != nil {
		tt.hooker.Hook(err)
	}
	return
}

type taskOptions struct {
	hooker Hooker[error]
}

func WithTaskErrorHooker(hooker Hooker[error]) Option[*taskOptions] {
	return FuncOption[*taskOptions](func(opts *taskOptions) {
		opts.hooker = hooker
	})
}

func NewLiteTask[T any](h func(T) error, v T, opts ...Option[*taskOptions]) Task {
	lt := &liteTask[T]{
		h: h,
		v: v,
	}
	dopts := &taskOptions{hooker: lt}
	for _, opt := range opts {
		opt.Apply(dopts)
	}
	lt.hooker = dopts.hooker
	return lt
}
