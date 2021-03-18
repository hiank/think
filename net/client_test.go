package net

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/pool"
	"google.golang.org/protobuf/proto"
	"gotest.tools/v3/assert"
)

type testClientFrame struct {
}

func (tf *testClientFrame) Dial(ctx context.Context, target string) (Conn, error) {
	return nil, errors.New("empty")
}

func (tf *testClientFrame) Handle(proto.Message) error {
	return nil
}

func TestNewClient(t *testing.T) {

	ctx, dialer, handler := context.Background(), &testClientFrame{}, &testClientFrame{}
	client := NewClient(ctx, dialer, handler)
	assert.Equal(t, client.ctx, ctx)
	assert.Equal(t, client.dialer, dialer)
	assert.Equal(t, client.recvHandler, handler)

	assert.Assert(t, client.hubPool != nil)
}

func TestClientAutoHub(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	framer := &testClientFrame{}
	client := NewClient(ctx, framer, framer)

	assert.Assert(t, client.autoHub("test") != nil)
	assert.Equal(t, client.autoHub("test"), client.autoHub("test"), "多次获取的hub需是同一个")

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

	loopRecv(ctx, ReciverFunc(func() (msg *pb.Message, err error) {
		msg, ok := <-msgChan
		if !ok {
			err = io.EOF
		}
		return
	}), pool.HandlerFunc(func(val proto.Message) error {

		msg := val.(*pb.Message)
		assert.Equal(t, msg.GetKey(), "lvws")
		uid = int(msg.GetSenderUid())
		return nil
	}))

	assert.Equal(t, uid, 128)

	cancel()
}

type testDialerFunc func(ctx context.Context, target string) (Conn, error)

func (td testDialerFunc) Dial(ctx context.Context, target string) (Conn, error) {
	return td(ctx, target)
}

type testConn struct {
	key string
	Sender
	Reciver
}

func (tc *testConn) Key() string {
	return tc.key
}

func (tc *testConn) Close() error {
	return nil
}

type testReciver struct {
	msgChan chan *pb.Message
}

func (tr *testReciver) Recv() (msg *pb.Message, err error) {

	msg, ok := <-tr.msgChan
	if !ok {
		err = io.EOF
	}
	return
}

func TestClientPush(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	outCh, sendCh := make(chan *pb.Message), make(chan *pb.Message)
	reciver := &testReciver{msgChan: make(chan *pb.Message)}

	client := NewClient(ctx, testDialerFunc(func(ctx context.Context, target string) (Conn, error) {

		return &testConn{
			key: target,
			Sender: SenderFunc(func(m *pb.Message) error {
				sendCh <- &pb.Message{Key: target, SenderUid: m.GetSenderUid()}
				return nil
			}),
			Reciver: reciver,
		}, nil
	}), pool.HandlerFunc(func(i proto.Message) error {
		outCh <- i.(*pb.Message)
		return nil
	}))

	msg := &pb.Message{Key: "testPush", SenderUid: 121}
	client.Push(msg)

	sendMsg := <-sendCh
	assert.Equal(t, sendMsg.GetKey(), "testPush")

	msg = &pb.Message{Key: "testRecv"}
	reciver.msgChan <- msg
	handleMsg := <-outCh

	assert.Equal(t, handleMsg, msg)

	cancel()
}
