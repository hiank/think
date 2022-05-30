package run_test

import (
	"container/list"
	"context"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/hiank/think/run"
)

type taskList struct {
	l    *list.List
	cond *sync.Cond
}

func newTaskList() *taskList {
	tl := &taskList{
		l:    list.New(),
		cond: sync.NewCond(&sync.Mutex{}),
	}
	return tl
}

func (tl *taskList) push(t run.Task) {
	tl.cond.L.Lock()
	defer tl.cond.L.Unlock()
	tl.l.PushBack(t)
	tl.cond.Signal()
}

func (tl *taskList) shift() (t run.Task) {
	tl.cond.L.Lock()
	defer tl.cond.L.Unlock()
	if tl.l.Len() == 0 {
		tl.cond.Wait()
	}
	if elm := tl.l.Front(); elm != nil {
		t = tl.l.Remove(elm).(run.Task)
	}
	return
}

func (tl *taskList) close() {
	tl.cond.L.Lock()
	defer tl.cond.L.Unlock()
	tl.cond.Signal()
}

type tasker struct {
	ctx context.Context
	tl  *taskList
	io.Closer
}

func NewTasker(ctx context.Context, timeout time.Duration) *tasker {
	ctx, cancel := context.WithCancel(ctx)
	healthy := run.NewHealthy()
	t := &tasker{
		ctx:    ctx,
		tl:     newTaskList(),
		Closer: run.NewHealthyCloser(healthy, cancel),
	}
	go healthy.Monitoring(ctx, t.tl.close)
	go t.looprocess(timeout)
	return t
}

func (t *tasker) tick(timeout time.Duration) (<-chan time.Time, func(time.Duration)) {
	if timeout > 0 {
		ticker := time.NewTicker(timeout)
		return ticker.C, ticker.Reset
	}
	return nil, func(d time.Duration) {}
}

func (t *tasker) looprocess(timeout time.Duration) {
	defer t.Close()
	tick, reset := t.tick(timeout)
	tt := time.Now().UnixMilli()
	for {
		tk := t.tl.shift()
		select {
		case <-t.ctx.Done(): //closed
			return
		case <-tick: //timeout
			return
		default:
			reset(timeout)
			tk.Process()
			ctt := time.Now().UnixMilli()
			fmt.Println("looprocess", ctt-tt)
			tt = ctt
		}
	}
}

func (t *tasker) Add(tk run.Task) (err error) {
	select {
	case <-t.ctx.Done():
		err = t.ctx.Err()
	default:
		t.tl.push(tk)
	}
	return err
}

// func (*tasker) internalOnly() {}

func BenchmarkTasker(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	tasker := run.NewTasker(ctx, time.Millisecond*200)
	max := 1000
	wait := new(sync.WaitGroup)
	wait.Add(max)
	for i := 0; i < max; i++ {
		err := tasker.Add(run.NewLiteTask(func(t int) error {
			<-time.After(time.Millisecond)
			// fmt.Println(t)
			wait.Done()
			return nil
		}, i+1, nil))
		if err != nil {
			b.Error(err)
		}
	}
	wait.Wait()
}

func BenchmarkTasker2(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	tasker := NewTasker(ctx, time.Millisecond*200)
	max := 1000
	wait := new(sync.WaitGroup)
	wait.Add(max)
	// ticker := time.NewTicker(time.Millisecond)
	for i := 0; i < max; i++ {
		err := tasker.Add(run.NewLiteTask(func(t int) error {
			ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
			defer cancel()
			tt := time.Now().UnixMilli()
			// <-time.After(time.Millisecond)
			// fmt.Println(t)
			// ticker.Reset(time.Millisecond)
			// <-time.Tick(time.Millisecond)
			// <-ticker.C
			<-ctx.Done()
			fmt.Println("task", time.Now().UnixMilli()-tt)
			wait.Done()
			return nil
		}, i+1, nil))
		if err != nil {
			b.Error(err)
		}
	}
	wait.Wait()
}

func BenchmarkSelectDefault(b *testing.B) {
	ctx, cnt := context.Background(), 10000000
	for i := 0; i < cnt; i++ {
		select {
		case <-ctx.Done():
		default:
		}
	}
}
