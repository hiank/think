package core_test

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"

	"gotest.tools/v3/assert"

	"github.com/hiank/think/core"
)

type testMessage struct {
	key string
}

func (tm *testMessage) GetKey() string {
	return tm.key
}

func (tm *testMessage) GetValue() *any.Any {
	return nil
}

type testConn struct {
	recvChan chan core.Message
	sendChan chan<- core.Message
	key      string
	once     sync.Once
}

func (tc *testConn) GetKey() string {
	return tc.key
}

func (tc *testConn) Recv() (core.Message, error) {

	if msg, ok := <-tc.recvChan; ok {
		return msg, nil
	}
	return nil, io.EOF
}

func (tc *testConn) Send(msg core.Message) error {

	tc.sendChan <- msg
	return nil
}

func (tc *testConn) Close() error {
	tc.once.Do(func() {
		close(tc.recvChan)
	})
	return nil
}

func TestPoolListen(t *testing.T) {

	pool, recvChan, sendChan := core.NewPool(context.Background()), make(chan core.Message), make(chan core.Message)
	conn := &testConn{
		recvChan: recvChan,
		sendChan: sendChan,
		key:      "test",
	}
	wait, notice := make(chan bool), make(chan bool)
	go func(wait chan bool) {
		close(wait)
		pool.Listen(conn, core.MessageHandlerTypeFunc(func(core.Message) error {
			time.Sleep(time.Millisecond * 100)
			notice <- true
			return nil
		}))
	}(wait)
	<-wait
	_, ok := pool.Get("test")
	assert.Assert(t, ok, "监听后，自动加入缓存")

	recvChan <- nil
	assert.Assert(t, <-notice, "接收端收到消息后，handler会被执行")
}

func TestPoolAddDelGetSafe(t *testing.T) {

	pool := core.NewPool(context.Background())
	for i := 0; i < 10000; i++ {
		switch rand.Intn(4) {
		case 0:
			fallthrough
		case 1:
			// go pool.Add(&testConn{key: fmt.Sprintf("test%d", rand.Intn(10))})
			go pool.Add(fmt.Sprintf("test%d", rand.Intn(10)), nil)
		case 2:
			go pool.Del(fmt.Sprintf("test%d", rand.Intn(10)))
		case 3:
			go pool.Get(fmt.Sprintf("test%d", rand.Intn(10)))
		}
	}
	assert.Assert(t, true, "上述多线程调用，如果没有panic，表明是线程安全的")
}

func TestPoolListenOverThenConnDeleted(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	pool, recvChan := core.NewPool(ctx), make(chan core.Message)
	conn := &testConn{
		recvChan: recvChan,
		key:      "test",
	}
	wait, errChan := make(chan bool), make(chan error)
	go func(wait chan bool) {
		close(wait)
		errChan <- pool.Listen(conn, core.MessageHandlerTypeFunc(func(core.Message) error {
			return nil
		}))
	}(wait)
	<-wait
	_, ok := pool.Get("test")
	assert.Assert(t, ok, "监听后，自动加入缓存")

	// recvChan <- nil
	// close(recvChan)
	cancel()
	recvChan <- nil

	// time.Sleep(time.Millisecond * 10)
	t.Log(<-errChan)
	_, ok = pool.Get("test")
	assert.Assert(t, !ok, "监听结束后，需要自动删除conn")
}

func TestPoolPush(t *testing.T) {

	pool := core.NewPool(context.Background())
	err := <-pool.Push(&testMessage{key: "test1"})
	assert.Assert(t, err != nil, "没有对应MessageHub时，Push 会返回错误")

	// wait := make(chan bool, 1)
	// go func() {
	// 	wait <- true
	// 	pool.Listen(nil, nil)
	// }()
	// <-wait
	// assert.Equal(t, (<-pool.Push(&testMessage{key: "test1"})).Error(), "test error", "处理结果必须正确返回给调用方")
}
