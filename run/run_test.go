package run_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hiank/think/run"
	"gotest.tools/v3/assert"
)

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
		err := tasker.Add(run.NewLiteTask(func(t int) error { return nil }, 1, nil))
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

	tasker := run.NewTasker(ctx, time.Second)
	pperr := make(chan error)
	err := errors.New("equal failed")
	tasker.Add(run.NewLiteTask(func(t int) error {
		if t != 10 {
			return err
		}
		return nil
	}, 11, pperr))

	err1 := <-pperr
	assert.Equal(t, err1, err)

	tasker.Add(run.NewLiteTask(func(t int) error {
		if t != 10 {
			return err
		}
		return nil
	}, 1, nil))
	//check run
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

}
