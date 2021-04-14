package rpc

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/pb"
	tg "github.com/hiank/think/net/rpc/pb"
	td "github.com/hiank/think/net/rpc/testdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

type testDonceHandler struct {
	postCh chan *pb.Message
}

func (tdh *testDonceHandler) HandleGet(msg *pb.Message) (rlt *pb.Message, err error) {
	return msg, nil
}

func (tdh *testDonceHandler) HandlePost(msg *pb.Message) error {
	tdh.postCh <- msg
	return nil
}

type testConn struct {
	key string
	net.Sender
	net.Reciver
}

func (tc *testConn) Key() string {
	return tc.key
}

func (tc *testConn) Close() error {
	return nil
}

func TestNewServeHelper(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	addr := ":10224"
	helper := NewServeHelper(ctx, ":10224", &testDonceHandler{})
	assert.Assert(t, helper.ctx != nil)
	assert.Equal(t, helper.addr, addr)
	assert.Assert(t, helper.connChan != nil)
	assert.Assert(t, helper.pipeServer != nil)
	cancel()
}

func TestServeHelperClose(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	helper := NewServeHelper(ctx, ":10224", &testDonceHandler{})
	go func() {
		helper.Close()
	}()
	_, ok := <-helper.connChan
	assert.Assert(t, !ok)
	select {
	case <-ctx.Done():
		assert.Assert(t, false, "关闭不能影响父context")
	default:
		assert.Assert(t, true, "关闭不能影响父context")
	}
	<-helper.ctx.Done()
	cancel()
}

func TestServeHelperAccept(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	helper := NewServeHelper(ctx, ":10224", &testDonceHandler{})

	t.Run("ok", func(t *testing.T) {

		go func() {
			helper.connChan <- &testConn{key: "OK"}
		}()
		conn, err := helper.Accept()
		assert.Assert(t, err == nil)
		assert.Equal(t, conn.Key(), "OK")
	})

	t.Run("closed", func(t *testing.T) {

		go func() {
			helper.Close()
		}()

		conn, err := helper.Accept()
		assert.Assert(t, conn == nil)
		assert.Equal(t, err, io.EOF)
	})

	cancel()
}

func TestServeHelperListenAndServe(t *testing.T) {

	postCh := make(chan *pb.Message, 1)
	ctx, cancel := context.WithCancel(context.Background())
	helper := NewServeHelper(ctx, ":10224", &testDonceHandler{postCh: postCh})
	closeChan := make(chan error)

	go func() {
		closeChan <- helper.ListenAndServe()
	}()

	t.Run("server can be dial", func(t *testing.T) {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		cc, err := grpc.DialContext(ctxWithTimeout, "localhost:10224", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name))) //NOTE: block 为阻塞直到ready，insecure 为不需要验证的
		assert.Assert(t, err == nil, err)
		err = cc.Close()
		assert.Assert(t, err == nil, err)

	})

	t.Run("get", func(t *testing.T) {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		cc, _ := grpc.DialContext(ctxWithTimeout, "localhost:10224", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name))) //NOTE: block 为阻塞直到ready，insecure 为不需要验证的
		// if err == nil {
		val, _ := anypb.New(&td.G_Example{})
		pipe, msg := tg.NewPipeClient(cc), &pb.Message{Key: "TestGet", Value: val}
		recvMsg, err := pipe.Donce(ctx, msg)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, recvMsg.GetKey(), msg.GetKey())
		cc.Close()
		// }
		// return
	})

	t.Run("post", func(t *testing.T) {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		cc, _ := grpc.DialContext(ctxWithTimeout, "localhost:10224", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name))) //NOTE: block 为阻塞直到ready，insecure 为不需要验证的
		val, _ := anypb.New(&td.P_Example{})
		pipe, msg := tg.NewPipeClient(cc), &pb.Message{Key: "TestPost", Value: val}
		recvMsg, err := pipe.Donce(ctx, msg)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, recvMsg.GetKey(), "", "Post返回的数据需为空数据，没有具体值")

		msg = <-postCh
		assert.Equal(t, msg.GetKey(), "TestPost", "需要检测到HandlePost的处理")
		cc.Close()
	})

	t.Run("donce: not support message", func(t *testing.T) {

		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		cc, _ := grpc.DialContext(ctxWithTimeout, "localhost:10224", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name))) //NOTE: block 为阻塞直到ready，insecure 为不需要验证的
		val, _ := anypb.New(&td.S_Example{})
		pipe, msg := tg.NewPipeClient(cc), &pb.Message{Key: "TestStream", Value: val}
		_, err := pipe.Donce(ctx, msg)
		assert.Assert(t, err != nil, "不支持的数据类型，需返回错误")
		cc.Close()
	})

	t.Run("stream", func(t *testing.T) {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		cc, _ := grpc.DialContext(ctxWithTimeout, "localhost:10224", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name))) //NOTE: block 为阻塞直到ready，insecure 为不需要验证的
		val, _ := anypb.New(&td.S_Example{})
		pipe := tg.NewPipeClient(cc)
		msg := &pb.Message{Key: "TestStream", Value: val}
		stream, err := pipe.Link(ctx)
		assert.Assert(t, err == nil, err)
		err = stream.Send(msg)
		assert.Assert(t, err == nil, err)
		conn, _ := helper.Accept()
		wait := make(chan bool)
		go func() {

			for {
				msg, err := conn.Recv()
				if err != nil {
					close(wait)
					return
				}
				conn.Send(msg)
			}
		}()

		value1, _ := anypb.New(&td.S_Example{Value: "1"})
		stream.Send(&pb.Message{Value: value1})

		recvMsg, err := stream.Recv()
		assert.Assert(t, err == nil, err)
		var tmpVal1 td.S_Example
		err = recvMsg.GetValue().UnmarshalTo(&tmpVal1)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, tmpVal1.GetValue(), "1")

		conn.Close()
		<-wait

		cc.Close()
	})

	t.Run("server close", func(t *testing.T) {

		helper.Close()
		err := <-closeChan
		assert.Assert(t, err == nil, err)
	})

	cancel()
}
