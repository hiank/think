package net_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/pb"
	"google.golang.org/protobuf/proto"
	"gotest.tools/v3/assert"
)

type testClientFrame struct {
}

func (tf *testClientFrame) Dial(ctx context.Context, target string) (net.Conn, error) {
	return nil, errors.New("empty")
}

func (tf *testClientFrame) Handle(proto.Message) error {
	return nil
}

func TestNewClient(t *testing.T) {

	ctx, dialer, handler := context.Background(), &testClientFrame{}, &testClientFrame{}
	client := net.NewClient(ctx, dialer, handler)
	assert.Equal(t, net.Export_getClientCtx(client), ctx)
	assert.Equal(t, net.Export_getClientDialer(client), dialer)
	assert.Equal(t, net.Export_getClientRecvHandler(client), handler)

	assert.Assert(t, net.Export_getClientHubPool(client) != nil)
}

func TestClientAutoHub(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	framer := &testClientFrame{}
	client := net.NewClient(ctx, framer, framer)

	autoHub := net.Export_ClientAutoHub(client)
	assert.Assert(t, autoHub("test") != nil)
	assert.Equal(t, autoHub("test"), autoHub("test"), "多次获取的hub需是同一个")

	time.Sleep(time.Second)
	cancel()
}

func TestLoopAccept(t *testing.T) {

}

func TestLoopRecv(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// tr := &testHandleRecover{msgChan: make(chan *pb.Message)}
	msgChan := make(chan *pb.Message)
	var uid int
	go func() {
		msgChan <- &pb.Message{Key: "lvws", SenderUid: 128}
		close(msgChan)
	}()

	net.Export_LoopRecv(ctx, net.ReciverFunc(func() (msg *pb.Message, err error) {
		msg, ok := <-msgChan
		if !ok {
			err = io.EOF
		}
		return
	}), net.HandlerFunc(func(msg *pb.Message) error {
		assert.Equal(t, msg.GetKey(), "lvws")
		uid = int(msg.GetSenderUid())
		return nil
	}))

	assert.Equal(t, uid, 128)

	cancel()
}

type testDialerFunc func(ctx context.Context, target string) (net.Conn, error)

func (td testDialerFunc) Dial(ctx context.Context, target string) (net.Conn, error) {
	return td(ctx, target)
}

type testConn struct {
	key string
	// net.Sender
	msgChan chan *pb.Message
	recvCh  chan *pb.Message
}

func newTestConn(key string) *testConn {
	return &testConn{
		key:     key,
		msgChan: make(chan *pb.Message, 1),
		recvCh:  make(chan *pb.Message),
	}
}

func newTestConnWithChan(key string, msgChan chan *pb.Message, recvCh chan *pb.Message) *testConn {
	return &testConn{key: key, msgChan: msgChan, recvCh: recvCh}
}

func (tc *testConn) Key() string {
	return tc.key
}

func (tc *testConn) Recv() (msg *pb.Message, err error) {
	msg, ok := <-tc.recvCh
	if !ok {
		err = io.EOF
	}
	return
}

func (tc *testConn) Send(msg *pb.Message) error {
	// tc.msg = msg
	tc.msgChan <- msg
	return nil
}

func (tc *testConn) Close() error {
	return nil
}

func TestClientPush(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	outCh, sendCh, recvCh := make(chan *pb.Message), make(chan *pb.Message), make(chan *pb.Message)

	client := net.NewClient(ctx, testDialerFunc(func(ctx context.Context, target string) (net.Conn, error) {

		return newTestConnWithChan(target, sendCh, recvCh), nil
	}), net.HandlerFunc(func(msg *pb.Message) error {
		outCh <- msg
		return nil
	}))

	msg := &pb.Message{Key: "testPush", SenderUid: 121}
	client.Push(msg)

	sendMsg := <-sendCh
	assert.Equal(t, sendMsg.GetKey(), "testPush")

	msg = &pb.Message{Key: "testRecv"}
	recvCh <- msg
	handleMsg := <-outCh

	assert.Equal(t, handleMsg, msg)

	cancel()
}
