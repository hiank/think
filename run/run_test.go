package run_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/hiank/think/run"
	"gotest.tools/v3/assert"
)

type tmpErrorHooker chan error

func (teh tmpErrorHooker) Hook(err error) {
	teh <- err
}

func TestHealthy(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	h := run.NewHealthy()
	e := make(chan bool)
	go h.Monitoring(ctx, func() {
		close(e)
	})
	cancel()
	<-e
	<-h.DoneContext().Done()

	h.Monitoring(context.Background(), func() {})
	t.Log("call again no response")
}

func TestTasker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	t.Run("timeout", func(t *testing.T) {
		tasker := run.NewTasker(ctx, time.Millisecond*10)
		<-time.After(time.Millisecond * 100)
		err := tasker.Add(run.NewLiteTask(func(t int) error { return nil }, 1))
		assert.Assert(t, err != nil, "closed because timeout")
	})

	t.Run("add to stopped tasker", func(t *testing.T) {
		tasker := run.NewTasker(ctx, time.Second)
		tasker.Close()
		//wait context canceled
		<-time.After(time.Microsecond * 100)
		err := tasker.Add(run.NewLiteTask(func(t int) error {
			return nil
		}, 1))
		assert.Assert(t, err != nil, "context canceled")
	})

	t.Run("ctx canceled", func(t *testing.T) {
		for i := 0; i < 10000; i++ {
			go func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				tasker := run.NewTasker(ctx, time.Millisecond*100)
				tasker.Add(run.NewLiteTask(func(t int) error {
					<-time.After(time.Millisecond)
					return nil
				}, 1))
			}(t)
		}
	})

	t.Run("stop immediately", func(t *testing.T) {
		tasker := run.NewTasker(ctx, time.Second)
		pp, wait := make(chan int, 10), make(chan bool)
		for i := 0; i < 10; i++ {
			err := tasker.Add(run.NewLiteTask(func(t int) error {
				if t == 0 {
					close(wait)
				}
				<-time.After(time.Millisecond * 10)
				pp <- t
				return nil
			}, i))
			assert.Equal(t, err, nil)
		}
		<-wait //wait for process goroutine started
		tasker.Close()
		assert.Equal(t, <-pp, 0)

		select {
		case <-pp:
			assert.Assert(t, false, "stoped")
		case <-time.After(time.Millisecond * 100):
			assert.Assert(t, true, "cannot process any task except first one")
		}
	})
	t.Run("loop add", func(t *testing.T) {
		tasker := run.NewTasker(ctx, time.Second)
		max := 100
		wait := new(sync.WaitGroup)
		wait.Add(max)
		for i := 0; i < max; i++ {
			go func(t *testing.T) {
				tasker.Add(run.NewLiteTask(func(t int) error {
					<-time.After(time.Millisecond)
					wait.Done()
					return nil
				}, 1))
			}(t)
		}
		wait.Wait()
	})
	t.Run("unrecoverable", func(t *testing.T) {
		tasker, c := run.NewTasker(ctx, time.Second), make(chan int, 10)
		tasker.Add(run.NewLiteTask(func(t int) error {
			c <- t
			return nil
		}, 1))
		tasker.Add(run.NewLiteTask(func(t int) error {
			c <- t
			return run.ErrUnrecoverable
		}, 2))
		tasker.Add(run.NewLiteTask(func(t int) error {
			c <- t
			return nil
		}, 3))
		<-time.After(time.Millisecond * 100)
		assert.Equal(t, len(c), 2)
		assert.Equal(t, <-c, 1)
		assert.Equal(t, <-c, 2)
		assert.Equal(t, len(c), 0)
	})
	// return
	tasker := run.NewTasker(ctx, time.Second)
	// pperr := make(chan error)
	hooker := make(tmpErrorHooker)
	err := errors.New("equal failed")
	tasker.Add(run.NewLiteTask(func(t int) error {
		if t != 10 {
			return err
		}
		return nil
	}, 11, run.WithTaskErrorHooker(hooker)))

	err1 := <-hooker
	assert.Equal(t, err1, err)

	tasker.Add(run.NewLiteTask(func(t int) error {
		if t != 10 {
			return err
		}
		return nil
	}, 1))
	//check run
}

func TestCloseChanWithCache(t *testing.T) {
	c, exit := make(chan int, 24), make(chan bool)
	go func(t *testing.T) {
		cnt := 0
		for range c {
			<-time.After(time.Millisecond * 10)
			cnt++
			// t.Log("do", cnt)
		}
		assert.Equal(t, cnt, 10, "chan will work after closed until cache cleaned")
		close(exit)
	}(t)
	for i := 0; i < 10; i++ {
		c <- i
	}
	close(c)
	// t.Log("after close")
	<-exit
}

func TestCustom(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	c := make(chan int)
	var outv int = 1
	go close(c)
L:
	for {
		select {
		case <-ctx.Done():
			break L
		case v, ok := <-c:
			outv = v
			if !ok {
				break L
			}
		}
	}
	assert.Equal(t, outv, 0)

	t.Run("range channel", func(t *testing.T) {
		c, v := make(chan int), 1
		// var v int = 1
		go close(c)
		for v = range c {
		}
		assert.Equal(t, v, 1, "when channel closed, range break without set value")
	})

	// var empty chan<- bool
	// empty <- true
}

func TestBitValue(t *testing.T) {
	var v8 int8 = -1
	for i := 0; i < 8; i++ {
		var tv int8 = 1
		tv <<= i
		assert.Equal(t, tv&v8, tv)
	}
	///
	var v1 int8 = 1
	v2 := (v1 << 6) - 1
	v2 |= (v1 << 6) | (v1 << 7)
	assert.Equal(t, v2, int8(-1))

	var v3 uint8
	v3 -= 1
	assert.Equal(t, v3, uint8(255))
}

// func TestMilliAfter(t *testing.T) {
// 	ticker := time.NewTicker(time.Millisecond)
// 	tt := time.Now().UnixMilli()
// 	// t.Log(time.Now().UnixMilli() - tt)
// 	for i := 0; i < 100; i++ {
// 		ticker.Reset(time.Millisecond)
// 		<-ticker.C
// 		ctt := time.Now().UnixMilli()
// 		// t.Log(ctt - tt)
// 		assert.Assert(t, ctt-tt > 1, "ticker响应会有10ms的额外开销")
// 		tt = ctt
// 	}

// }
