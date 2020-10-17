package ws_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/gorilla/websocket"
	"github.com/hiank/think/core"
	"github.com/hiank/think/core/pb"
	"github.com/hiank/think/core/ws"
	td "github.com/hiank/think/core/ws/testdata"
	"google.golang.org/protobuf/proto"
	"gotest.tools/assert"
)

func dail(addr string, token string) (*websocket.Conn, *http.Response, error) {

	u := url.URL{Scheme: "ws", Host: addr + ":8022", Path: "/ws"}
	return websocket.DefaultDialer.Dial(u.String(), http.Header{"token": {token}})
}

func startAutoServe(handle func()) {

	ctx, cancel := context.WithCancel(context.Background())
	notice := make(chan byte, 1)
	go func() {
		notice <- 0
		ws.ListenAndServe(ctx, "localhost", core.MessageHandlerTypeFunc(func(msg core.Message) error {

			req := &td.Request{Value: msg.GetKey()}
			anyMsg, _ := ptypes.MarshalAny(req)

			var writer ws.Writer
			return writer.Handle(&pb.Message{Key: msg.GetKey(), Value: anyMsg})
		}))
		close(notice)
	}()
	<-notice

	handle()

	cancel()
	<-notice
}

func TestServerStart(t *testing.T) {

	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		errChan <- ws.ListenAndServe(ctx, "localhost", nil)
		errChan <- ws.ListenAndServe(ctx, "localhost", nil) //NOTE: 测试ws服务可以立即重新开启，方便测试构建
	}()
	cancel()
	err := <-errChan
	assert.Assert(t, err != nil, err)
	cancel()
	err = <-errChan
	assert.Assert(t, err != nil, err)
}

func TestConnect(t *testing.T) {

	startAutoServe(func() {

		conn, _, err := dail("localhost", "token1")
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()

		anyMsg, _ := ptypes.MarshalAny(&td.Request{Value: "Ivws"})
		buf, err := proto.Marshal(anyMsg)
		if err != nil {
			t.Error(err)
		}
		if err = conn.WriteMessage(websocket.BinaryMessage, buf); err != nil {
			t.Error(err)
		}
		if _, buf, err = conn.ReadMessage(); err != nil {
			t.Error(err)
		}
		var anyMsgRecv = new(any.Any)
		if err = proto.Unmarshal(buf, anyMsgRecv); err != nil {
			t.Error(err)
		}

		var resp = &td.Request{}
		ptypes.UnmarshalAny(anyMsgRecv, resp)
		assert.Equal(t, resp.GetValue(), "token1", "需要收到期望的数据")
	})
}

func TestConnectHeavy(t *testing.T) {

	var test = func(t *testing.T, token string) {

		conn, _, err := dail("localhost", token)
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()

		anyMsg, _ := ptypes.MarshalAny(&td.Request{Value: "Ivws"})
		buf, err := proto.Marshal(anyMsg)
		if err != nil {
			t.Error(err)
		}
		if err = conn.WriteMessage(websocket.BinaryMessage, buf); err != nil {
			t.Error(err)
		}
		if _, buf, err = conn.ReadMessage(); err != nil {
			t.Error(err)
		}
		var anyMsgRecv = new(any.Any)
		if err = proto.Unmarshal(buf, anyMsgRecv); err != nil {
			t.Error(err)
		}

		var resp = &td.Request{}
		ptypes.UnmarshalAny(anyMsgRecv, resp)
		assert.Equal(t, resp.GetValue(), token, "需要收到期望的数据")
	}

	startAutoServe(func() {

		max, wait := 100, new(sync.WaitGroup) //NOTE: 本地端口数限制，需要寻找其它方式模拟超级并发
		wait.Add(max)
		for i := 0; i < max; i++ {
			go func(t *testing.T, i int) {
				test(t, fmt.Sprintf("token%d", i))
				wait.Done()
			}(t, i)
		}
		wait.Wait()
	})
}
