package run

import (
	"container/list"
	"context"
	"sync"
	"time"

	"k8s.io/klog/v2"
)

type tasker struct {
	ctx     context.Context
	cancel  context.CancelFunc
	tc      chan *Task
	once    *sync.Once
	timeout time.Duration
}

func NewTasker(ctx context.Context, timeout time.Duration) Tasker {
	ctx, cancel := context.WithCancel(ctx)
	return &tasker{
		ctx:     ctx,
		cancel:  cancel,
		tc:      make(chan *Task, 16),
		once:    new(sync.Once),
		timeout: timeout,
	}
}

func (t *tasker) process(c chan *Task) {
	for tk := range c {
		if err := tk.H(tk.V); err != nil && tk.C != nil {
			select {
			case tk.C <- err:
			case <-time.NewTicker(time.Second * 5).C:
				klog.Warning("cannot send error for invalid 'handle' in long try")
			}
		}
	}
}

func (t *tasker) loop() {
	l, ticker := list.New(), time.NewTicker(t.timeout)
	var pp chan *Task
	var wt *Task
L:
	for {
		select {
		case <-t.ctx.Done():
			if pp != nil {
				close(pp) //stop process goroutine
			}
			return
		case v := <-t.tc:
			if pp == nil {
				pp, wt = make(chan *Task), v
				go t.process(pp)
			} else {
				l.PushBack(v)
			}
		case pp <- wt:
			if em := l.Front(); em == nil {
				close(pp)
				pp, wt = nil, nil
			} else {
				wt = em.Value.(*Task)
			}
		case <-ticker.C:
			if pp == nil {
				break L
			}
		}
		ticker.Reset(t.timeout)
	}
	t.once = new(sync.Once)
}

func (tkr *tasker) Add(tk Task) (err error) {
	if err = tkr.ctx.Err(); err == nil {
		tkr.once.Do(func() {
			go tkr.loop()
		})
		tkr.tc <- &tk
	}
	return err
}

func (tkr *tasker) Stop() {
	tkr.cancel()
}
