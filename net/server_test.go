package net_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/set/codes"
	"google.golang.org/protobuf/proto"
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
	recvHandler := pool.HandlerFunc(func(m proto.Message) error { return err })
	net.Export_ListenOptionApply(net.WithRecvHandler(recvHandler))(dopts)
	assert.Equal(t, net.Export_getListenOptionsRecvHandler(dopts).Handle(nil), recvHandler.Handle(nil), "设置recvHandler后，opts中的值变为传入的handler")
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
			}), net.WithRecvHandler(pool.HandlerFunc(func(proto.Message) error {
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
			waitErr <- net.ListenAndServe(helper, net.WithRecvHandler(pool.HandlerFunc(func(proto.Message) error {
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
		err := net.ListenAndServe(newTestServeHelper(), net.WithContext(ctx), net.WithRecvHandler(pool.HandlerFunc(func(proto.Message) error {
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
			waitErr <- net.ListenAndServe(helper, net.WithRecvHandler(pool.HandlerFunc(func(m proto.Message) error { return nil })))
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

// func Test

// func TestNewServer(t *testing.T) {

// 	ctx, cancel := context.WithCancel(context.Background())
// 	helper := &testServeHelper{}
// 	msgHandler := pool.HandlerFunc(func(i proto.Message) error {
// 		return nil
// 	})
// 	srv := NewServer(ctx, helper, msgHandler)

// 	// assert.Equal(t, srv.ctx, ctx)
// 	assert.Equal(t, srv.helper, helper)
// 	assert.Assert(t, srv.hubPool != nil)
// 	assert.Assert(t, srv.recvHandler != nil)

// 	cancel()
// 	<-srv.ctx.Done()
// }

// func TestServerLoopAccept(t *testing.T) {

// 	ctx, cancel := context.WithCancel(context.Background())
// 	connCh := make(chan Conn)
// 	helper := &testServeHelper{
// 		ctx:    ctx,
// 		connCh: connCh,
// 	}
// 	srv := NewServer(ctx, helper, nil)

// 	go func() {
// 		connCh <- nil
// 		time.Sleep(time.Second)
// 		// cancel()
// 		close(connCh)
// 	}()
// 	// srv.loopAccept()
// 	loopAccept(srv.ctx, srv.helper.(Accepter), srv.handleAccept)

// 	cancel()
// }

// func TestServerListenAndServe(t *testing.T) {

// 	t.Run("ListenAndServe: context cancel", func(t *testing.T) {

// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 		srv := NewServer(
// 			ctx,
// 			&testServeHelper{
// 				ctx:    ctx,
// 				connCh: make(chan Conn),
// 			},
// 			nil)

// 		err := srv.ListenAndServe()
// 		assert.Assert(t, err != nil)

// 		cancel()
// 	})

// 	t.Run("ListenAndServe: hubPool response", func(t *testing.T) {

// 		connCh := make(chan Conn)

// 		ctx, cancel := context.WithCancel(context.Background())
// 		srv := NewServer(
// 			ctx,
// 			&testServeHelper{
// 				ctx:    ctx,
// 				connCh: connCh,
// 			},
// 			nil,
// 		)
// 		go func(t *testing.T) {
// 			connCh <- &testConn{
// 				key: "key1",
// 				Reciver: ReciverFunc(func() (*pb.Message, error) {
// 					<-ctx.Done()
// 					return nil, io.EOF
// 				}),
// 			}
// 			time.Sleep(time.Second)
// 			assert.Assert(t, srv.hubPool.GetHub("key1") != nil)

// 			cancel()
// 		}(t)
// 		srv.ListenAndServe()
// 	})
// }

// // type test

// func TestServerSend(t *testing.T) {

// 	// var value int
// 	connCh := make(chan Conn)
// 	ctx, cancel := context.WithCancel(context.Background())

// 	srv := NewServer(
// 		ctx,
// 		&testServeHelper{
// 			connCh: connCh,
// 			ctx:    ctx,
// 		},
// 		nil,
// 	)

// 	go srv.ListenAndServe()

// 	t.Run("Send: nil message", func(t *testing.T) {

// 		err := srv.Send(nil)
// 		assert.Equal(t, err, codes.Error(codes.ErrorNilValue), "发送nil消息，返回指定错误")
// 	})

// 	t.Run("Send: not existed hub for message", func(t *testing.T) {

// 		err := srv.Send(&pb.Message{Key: "notExistedKey"})
// 		assert.Equal(t, err, codes.Error(codes.ErrorWorkerNotExisted), "发送的消息无效的Key，返回指定错误")
// 	})

// 	t.Run("Send: success", func(t *testing.T) {

// 		sendCh := make(chan *pb.Message)

// 		connCh <- &testConn{
// 			key: "TestKey",
// 			Reciver: ReciverFunc(func() (*pb.Message, error) {
// 				<-ctx.Done()
// 				return nil, io.EOF
// 			}),
// 			Sender: SenderFunc(func(msg *pb.Message) error {
// 				sendCh <- msg
// 				return nil
// 			}),
// 		}

// 		time.Sleep(time.Second)

// 		err := srv.Send(&pb.Message{Key: "TestKey"})
// 		assert.Assert(t, err == nil)

// 		msg := <-sendCh
// 		assert.Equal(t, msg.GetKey(), "TestKey")

// 		time.Sleep(time.Second)
// 	})

// 	cancel()
// }

// func TestServerRecv(t *testing.T) {

// 	// var value int
// 	handleCh := make(chan interface{})
// 	connCh := make(chan Conn)
// 	ctx, cancel := context.WithCancel(context.Background())

// 	srv := NewServer(
// 		ctx,
// 		&testServeHelper{
// 			connCh: connCh,
// 			ctx:    ctx,
// 		},
// 		pool.HandlerFunc(func(i proto.Message) error {
// 			handleCh <- i
// 			return nil
// 		}),
// 	)

// 	go srv.ListenAndServe()

// 	recvCh := make(chan *pb.Message)
// 	connCh <- &testConn{
// 		key: "TestKey",
// 		Reciver: ReciverFunc(func() (*pb.Message, error) {
// 			return <-recvCh, nil
// 		}),
// 	}
// 	var msg *pb.Message
// 	recvCh <- msg
// 	assert.Equal(t, <-handleCh, msg)

// 	cancel()
// }

// func TestServerClose(t *testing.T) {

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	helper := &testServeHelper{
// 		ctx: ctx,
// 	}

// 	srv := NewServer(
// 		ctx,
// 		helper,
// 		nil,
// 	)

// 	helper.ctx = srv.ctx

// 	errChan, noticeChan := make(chan error), make(chan bool, 1)
// 	go func() {
// 		noticeChan <- true
// 		err := srv.ListenAndServe()
// 		errChan <- err
// 	}()

// 	<-noticeChan
// 	time.Sleep(time.Microsecond)

// 	srv.Close()

// 	err := <-errChan
// 	assert.Equal(t, err, io.EOF)
// }
