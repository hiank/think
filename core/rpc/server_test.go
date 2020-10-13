package rpc_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hiank/think/core"
	"github.com/hiank/think/core/pb"
	"github.com/hiank/think/core/rpc"
)

func TestServerStart(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	go func(t *testing.T) {
		time.Sleep(time.Second * 1)
		cancel()
	}(t)
	rpc.ListenAndServe(ctx, "localhost", nil)
	// rpc.ListenAndServe(ctx, "localhost", nil)
}

func startOneServer(ctx context.Context, wait *sync.WaitGroup, handler rpc.ReadHandler) {

	rpc.ListenAndServe(ctx, "localhost", handler)
	wait.Done()
}

type testReadHandler struct {
	rpc.IgnorePost
	// rpc.IgnoreStream
	rpc.IgnoreGet
}

func (th *testReadHandler) Handle(msg core.Message) (err error) {
	fmt.Println("Handle : ", msg.GetKey())
	return
}

func TestServerWrite(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	wait := new(sync.WaitGroup)
	wait.Add(1)
	go func() {
		// notice <- true
		startOneServer(ctx, wait, new(testReadHandler))
	}()
	time.Sleep(time.Second)
	t.Log(new(rpc.Writer).Handle(&pb.Message{
		Key:   "test",
		Value: nil,
	}))
	// cancel()
	
	cancel()
	wait.Wait()
}
