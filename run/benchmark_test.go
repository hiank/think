package run_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/hiank/think/run"
)

func BenchmarkTasker(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	tasker := run.NewTasker(ctx, time.Millisecond*20)
	max := 100
	wait := new(sync.WaitGroup)
	wait.Add(max)
	for i := 0; i < max; i++ {
		err := tasker.Add(run.NewLiteTask(func(t int) error {
			<-time.After(time.Millisecond * 1)
			// <-time.NewTicker(time.Millisecond * 10).C
			wait.Done()
			return nil
		}, i+1, nil))
		if err != nil {
			b.Error(err)
		}
	}
	wait.Wait()
}
