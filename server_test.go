package think_test

import (
	"testing"

	"github.com/hiank/think"

	"github.com/hiank/think/core/rpc"

	"gotest.tools/assert"
)

// func TestCloseWS(t *testing.T) {

// 	exit := make(chan error)

// 	go func() {
// 		exit <- think.ServeWSDefault("127.0.0.1")
// 	}()

// 	go func() {
// 		<-time.After(time.Second)
// 		// token.BackgroundLife().Kill()
// 	}()
// 	err := <-exit
// 	assert.Equal(t, err, http.ErrServerClosed)
// }

type testK8sHandler struct {
	rpc.IgnoreStream
	rpc.IgnoreGet
	rpc.IgnorePost
}

// func TestCloseK8s(t *testing.T) {

// 	exit := make(chan error)

// 	go func() {
// 		exit <- think.ServeRPC("127.0.0.1", &testK8sHandler{})
// 	}()

// 	go func() {
// 		<-time.After(time.Second)
// 		token.BackgroundLife().Kill()
// 	}()
// 	err := <-exit
// 	assert.Equal(t, err, nil)
// }

func TestTryMQ(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
		}
	}()
	think.TryMQ()
	assert.Assert(t, false, "目前无法连接到消息中间件，不能执行这个方法")
}

//testHandleMQ 这个用于测试 recover中对返回值赋值 调用者是可以收到的
func testHandleMQ() (err error) {

	defer func() {
		if recv := recover(); recv != nil {
			err = recv.(error)
		}
	}()
	think.TryMQ()
	return
}

func TestHandleMQ(t *testing.T) {

	err := testHandleMQ()
	t.Log(err)
}
