package net_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	testdatapb "github.com/hiank/think/net/testdata"

	"github.com/gorilla/websocket"
	"github.com/hiank/think/net"
	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/net/ws"
	"github.com/hiank/think/pool"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

type testAuther uint64

func (ta testAuther) Auth(token string) (uint64, error) {
	return uint64(ta), nil
}

func StartSimpleWSSvc(ctx context.Context) <-chan bool {
	srv := net.NewServer(ctx, ws.NewServeHelper(":8022", testAuther(1024)), pb.LiteHandler)
	// client := net.NewClient(context.Background(), rpc.Dialer, pool.HandlerFunc(func(m proto.Message) error {
	// 	return srv.Send(m.(*pb.Message))
	// }))
	pb.LiteHandler.DefaultHandler = pool.HandlerFunc(func(m proto.Message) error { //NOTE: 默认使用grpc 将消息转发到k8s集群中
		// client.Push(m.(*pb.Message))
		pbMsg := m.(*pb.Message)
		fmt.Printf("on handle: %v\n", pbMsg)
		srv.Send(&pb.Message{Key: strconv.FormatUint(pbMsg.GetSenderUid(), 10), Value: pbMsg.GetValue()})
		return nil
	})

	out := make(chan bool)
	go func() {
		err := srv.ListenAndServe()
		fmt.Print(err)
		out <- true
	}()
	return out
}

func TestDail(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	out := StartSimpleWSSvc(ctx)

	addr := "127.0.0.1:8022"

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	_, res, _ := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.Equal(t, res.StatusCode, http.StatusNonAuthoritativeInfo, "url需要包含Token")

	_, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"token": {"love-ws"}})
	assert.Assert(t, err == nil, err)

	cancel()
	<-out
}

func TestReadWrite(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	out := StartSimpleWSSvc(ctx)

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8022", Path: "/ws"}
	wsConn, _, _ := websocket.DefaultDialer.Dial(u.String(), http.Header{"token": {"love-ws"}})

	anyMsg, _ := anypb.New(&testdatapb.Test1{Name: "lovews"})
	buff, _ := proto.Marshal(anyMsg)
	err := wsConn.WriteMessage(websocket.BinaryMessage, buff)
	assert.Assert(t, err == nil, err)

	_, buff, err = wsConn.ReadMessage()
	assert.Assert(t, err == nil, err)

	var anyRlt anypb.Any
	proto.Unmarshal(buff, &anyRlt)

	msg, _ := anyRlt.UnmarshalNew()
	rltMsg := msg.(*testdatapb.Test1)
	assert.Equal(t, rltMsg.GetName(), "lovews")

	cancel()
	<-out
}
