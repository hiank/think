package net_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/hiank/think/utils/robust"

	"github.com/hiank/think/net/rpc"
	"github.com/hiank/think/token"

	"github.com/hiank/think/net"
	"gotest.tools/assert"
)

func TestCloseWS(t *testing.T) {

	exit := make(chan error)

	go func() {
		exit <- net.ServeWSDefault("127.0.0.1")
	}()

	go func() {
		<-time.After(time.Second)
		token.BackgroundLife().Kill()
	}()
	err := <-exit
	assert.Equal(t, err, http.ErrServerClosed)
}

type testK8sHandler struct {
	rpc.IgnoreStream
	rpc.IgnoreGet
	rpc.IgnorePost
}

func TestCloseK8s(t *testing.T) {

	exit := make(chan error)

	go func() {
		exit <- net.ServeRPC("127.0.0.1", &testK8sHandler{})
	}()

	go func() {
		<-time.After(time.Second)
		token.BackgroundLife().Kill()
	}()
	err := <-exit
	assert.Equal(t, err, nil)
}

func TestTryMQ(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
		}
	}()
	net.TryMQ()
	assert.Assert(t, false, "目前无法连接到消息中间件，不能执行这个方法")
}

//testHandleMQ 这个用于测试 recover中对返回值赋值 调用者是可以收到的
func testHandleMQ() (err error) {

	defer robust.Recover(robust.Warning, robust.ErrorHandle(func(e interface{}) {
		err = e.(error)
	}))
	net.TryMQ()
	return
}

func TestHandleMQ(t *testing.T) {

	err := testHandleMQ()
	t.Log(err)
}
