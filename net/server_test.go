package net_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/set/codes"
	"gotest.tools/v3/assert"
)

func TestLitenOption(t *testing.T) {
	dopts := net.Export_newDefaultListenOptions()
	dctx := net.Export_getListenOptionsCtx(dopts)
	assert.Assert(t, dctx != nil, "默认的context为context.Background()")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	net.Export_ListenOptionApply(net.WithContext(ctx))(dopts)
	assert.Assert(t, ctx != dctx)
	assert.Equal(t, ctx, net.Export_getListenOptionsCtx(dopts), "设置context后，opts中的值变为传入的context")

	assert.Assert(t, net.Export_getListenOptionsConnHandler(dopts) == nil, "默认的connHandler为nil")
	connHandler := func(net.Conn) {

	}
	net.Export_ListenOptionApply(net.WithConnHandler(connHandler))(dopts)
	assert.Assert(t, net.Export_getListenOptionsConnHandler(dopts) != nil, "设置connHandler后，opts中的值变为传入的handler")
	// net.Export_getListenOptionsConnHandler()

	err := errors.New("tmp")
	recvHandler := net.HandlerFunc(func(*pb.Message) error { return err })
	net.Export_ListenOptionApply(net.WithRecvHandler(recvHandler))(dopts)
	assert.Equal(t, net.Export_getListenOptionsRecvHandler(dopts).Handle(nil).Error(), recvHandler.Handle(nil).Error(), "设置recvHandler后，opts中的值变为传入的handler")
}

type testServeHelper struct {
	ctx    context.Context
	cancel context.CancelFunc
	connCh chan net.Conn
}

func newTestServeHelper() *testServeHelper {
	ctx, cancel := context.WithCancel(context.Background())
	return &testServeHelper{
		ctx:    ctx,
		cancel: cancel,
		connCh: make(chan net.Conn),
	}
}

func (th *testServeHelper) ListenAndServe() error {

	<-th.ctx.Done()
	return io.EOF
}

func (th *testServeHelper) Accept() (conn net.Conn, err error) {

	conn, ok := <-th.connCh
	if !ok {
		err = io.EOF
	}
	return
}

func (th *testServeHelper) Close() error {
	th.cancel()
	return nil
}

// type testConn struct {
// 	key string
// 	// net.Sender
// 	msgChan chan *pb.Message
// }

// func newTestConn(key string) *testConn {
// 	return &testConn{
// 		key:     key,
// 		msgChan: make(chan *pb.Message, 1),
// 	}
// }

// func (tc *testConn) Key() string {
// 	return tc.key
// }

// func (tc *testConn) Recv() (*pb.Message, error) {
// 	return nil, nil
// }

// func (tc *testConn) Send(msg *pb.Message) error {
// 	// tc.msg = msg
// 	tc.msgChan <- msg
// 	return nil
// }

func TestListenAndServe(t *testing.T) {
	t.Run("NoHelper-ReturnError", func(t *testing.T) {
		err := net.ListenAndServe(nil)
		assert.Equal(t, err, codes.Error(codes.ErrorNonHelper))
	})
	t.Run("NoConnHandlerAndNoRecvHandler-ReturnError", func(t *testing.T) {
		err := net.ListenAndServe(&testServeHelper{})
		assert.Equal(t, err, codes.Error(codes.ErrorNeedOneofConnRecvHandler))
	})
	t.Run("ExistConnHandler-RecvHandlerNoUseful", func(t *testing.T) {
		helper, waitErr := newTestServeHelper(), make(chan error)
		// var connNum, recvNum int
		numChan := make(chan int)
		go func() {
			waitErr <- net.ListenAndServe(helper, net.WithConnHandler(func(net.Conn) {
				numChan <- 1
			}), net.WithRecvHandler(net.HandlerFunc(func(*pb.Message) error {
				numChan <- 2
				return nil
			})))
		}()

		for i := 0; i < 10000; i++ {
			go func() {
				helper.connCh <- &testConn{}
			}()
		}

		ticker := time.NewTicker(time.Millisecond * 100)
	L:
		for {
			select {
			case num := <-numChan:
				assert.Equal(t, num, 1, "只能收到ConnHandler发送的消息")
			case <-ticker.C: //NOTE: 限定100ms
				helper.Close()
				break L
			}
		}
		// t.Log(<-waitErr)
		err := <-waitErr
		assert.Equal(t, err, io.EOF, err) //NOTE:testServeHelper 关闭返回的是io.EOF
	})
	t.Run("UseRecvHandler", func(t *testing.T) {
		helper, waitErr := newTestServeHelper(), make(chan error)
		// var connNum, recvNum int
		numChan := make(chan int)
		go func() {
			waitErr <- net.ListenAndServe(helper, net.WithRecvHandler(net.HandlerFunc(func(*pb.Message) error {
				numChan <- 2
				return nil
			})))
		}()

		for i := 0; i < 10000; i++ {
			go func() {
				helper.connCh <- &testConn{}
			}()
		}
		ticker := time.NewTicker(time.Millisecond * 100)
	L:
		for {
			select {
			case num := <-numChan:
				assert.Equal(t, num, 2, "没有ConnHandler情况下，RecvHandler需要响应")
			case <-ticker.C: //NOTE: 限定100ms
				helper.Close()
				break L
			}
		}
		err := <-waitErr
		assert.Equal(t, err, io.EOF, err) //NOTE:
	})
	t.Run("WithContext", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go cancel()
		err := net.ListenAndServe(newTestServeHelper(), net.WithContext(ctx), net.WithRecvHandler(net.HandlerFunc(func(*pb.Message) error {
			return nil
		})))
		assert.Equal(t, err, io.EOF, err) //NOTE:
	})
}

func TestLiteSender(t *testing.T) {
	t.Run("UseConnHandler", func(t *testing.T) {
		helper, waitErr := newTestServeHelper(), make(chan error)
		go func() {
			waitErr <- net.ListenAndServe(helper, net.WithConnHandler(func(net.Conn) {
			}))
		}()

		err := net.LiteSender.Send(nil)
		assert.Equal(t, err, codes.Error(codes.ErrorNonSupportLiteServe))
		helper.Close()
		<-waitErr
	})

	t.Run("UseRecvHandler", func(t *testing.T) {
		helper, waitErr := newTestServeHelper(), make(chan error)
		go func() {
			waitErr <- net.ListenAndServe(helper, net.WithRecvHandler(net.HandlerFunc(func(*pb.Message) error { return nil })))
		}()

		<-time.NewTicker(time.Millisecond * 10).C //NOTE: wait 10ms, for serve start
		err := net.LiteSender.Send(nil)
		assert.Equal(t, err, codes.Error(codes.ErrorNilValue), "try send nil message will get ErrorNilValue error")

		err = net.LiteSender.Send(&pb.Message{})
		assert.Equal(t, err, codes.Error(codes.ErrorWorkerNotExisted), err)

		conn := newTestConn("1024")
		helper.connCh <- conn //&testConn{key: "1024"}

		<-time.NewTicker(time.Millisecond * 10).C //NOTE: wait 10ms, for conn Accepted

		msg := &pb.Message{Key: "1024"}
		err = net.LiteSender.Send(msg)
		assert.Equal(t, err, nil, err)
		tmsg := <-conn.msgChan
		assert.Equal(t, msg, tmsg)

		helper.Close()
		<-waitErr
	})
}
