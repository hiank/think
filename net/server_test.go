package net

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/set/codes"
	"google.golang.org/protobuf/proto"
	"gotest.tools/v3/assert"
)

type testAccepter interface {
	Accept() (Conn, error)
}

type testAccepterFunc func() (Conn, error)

func (tf testAccepterFunc) Accept() (Conn, error) {
	return tf()
}

func TestNewServer(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	helper := &testServerHelper{}
	msgHandler := pool.HandlerFunc(func(i proto.Message) error {
		return nil
	})
	srv := NewServer(ctx, helper, msgHandler)

	// assert.Equal(t, srv.ctx, ctx)
	assert.Equal(t, srv.helper, helper)
	assert.Assert(t, srv.hubPool != nil)
	assert.Assert(t, srv.recvHandler != nil)

	cancel()
	<-srv.ctx.Done()
}

func TestServerLoopAccept(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	connCh := make(chan Conn)
	helper := &testServerHelper{
		ctx:    ctx,
		connCh: connCh,
	}
	srv := NewServer(ctx, helper, nil)

	go func() {
		connCh <- nil
		time.Sleep(time.Second)
		// cancel()
		close(connCh)
	}()
	srv.loopAccept()

	cancel()
}

type testServerHelper struct {
	ctx    context.Context
	connCh chan Conn
}

func (th *testServerHelper) ListenAndServe() error {

	<-th.ctx.Done()
	return io.EOF
}

func (th *testServerHelper) Accept() (conn Conn, err error) {

	conn, ok := <-th.connCh
	if !ok {
		err = io.EOF
	}
	return
}

func (th *testServerHelper) Close() error {
	return nil
}

func TestServerListenAndServe(t *testing.T) {

	t.Run("ListenAndServe: context cancel", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		srv := NewServer(
			ctx,
			&testServerHelper{
				ctx:    ctx,
				connCh: make(chan Conn),
			},
			nil)

		err := srv.ListenAndServe()
		assert.Assert(t, err != nil)

		cancel()
	})

	t.Run("ListenAndServe: hubPool response", func(t *testing.T) {

		connCh := make(chan Conn)

		ctx, cancel := context.WithCancel(context.Background())
		srv := NewServer(
			ctx,
			&testServerHelper{
				ctx:    ctx,
				connCh: connCh,
			},
			nil,
		)
		go func(t *testing.T) {
			connCh <- &testConn{
				key: "key1",
				Reciver: ReciverFunc(func() (*pb.Message, error) {
					<-ctx.Done()
					return nil, io.EOF
				}),
			}
			time.Sleep(time.Second)
			assert.Assert(t, srv.hubPool.GetHub("key1") != nil)

			cancel()
		}(t)
		srv.ListenAndServe()
	})
}

// type test

func TestServerSend(t *testing.T) {

	// var value int
	connCh := make(chan Conn)
	ctx, cancel := context.WithCancel(context.Background())

	srv := NewServer(
		ctx,
		&testServerHelper{
			connCh: connCh,
			ctx:    ctx,
		},
		nil,
	)

	go srv.ListenAndServe()

	t.Run("Send: nil message", func(t *testing.T) {

		err := srv.Send(nil)
		assert.Equal(t, err, codes.ErrorNilValue, "发送nil消息，返回指定错误")
	})

	t.Run("Send: not existed hub for message", func(t *testing.T) {

		err := srv.Send(&pb.Message{Key: "notExistedKey"})
		assert.Equal(t, err, codes.ErrorNotExisted, "发送的消息无效的Key，返回指定错误")
	})

	t.Run("Send: success", func(t *testing.T) {

		sendCh := make(chan *pb.Message)

		connCh <- &testConn{
			key: "TestKey",
			Reciver: ReciverFunc(func() (*pb.Message, error) {
				<-ctx.Done()
				return nil, io.EOF
			}),
			Sender: SenderFunc(func(msg *pb.Message) error {
				sendCh <- msg
				return nil
			}),
		}

		time.Sleep(time.Second)

		err := srv.Send(&pb.Message{Key: "TestKey"})
		assert.Assert(t, err == nil)

		msg := <-sendCh
		assert.Equal(t, msg.GetKey(), "TestKey")

		time.Sleep(time.Second)
	})

	cancel()
}

func TestServerRecv(t *testing.T) {

	// var value int
	handleCh := make(chan interface{})
	connCh := make(chan Conn)
	ctx, cancel := context.WithCancel(context.Background())

	srv := NewServer(
		ctx,
		&testServerHelper{
			connCh: connCh,
			ctx:    ctx,
		},
		pool.HandlerFunc(func(i proto.Message) error {
			handleCh <- i
			return nil
		}),
	)

	go srv.ListenAndServe()

	recvCh := make(chan *pb.Message)
	connCh <- &testConn{
		key: "TestKey",
		Reciver: ReciverFunc(func() (*pb.Message, error) {
			return <-recvCh, nil
		}),
	}
	var msg *pb.Message
	recvCh <- msg
	assert.Equal(t, <-handleCh, msg)

	cancel()
}
