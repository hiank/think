package rpc

import (
	"context"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	net1 "github.com/hiank/think/net"
	"github.com/hiank/think/net/pb"
	tg "github.com/hiank/think/net/rpc/pb"
	td "github.com/hiank/think/net/testdata"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

func startGrpcServer(ctx context.Context, t *testing.T, handlers ...func(*grpc.Server)) {

	lis, err := new(net.ListenConfig).Listen(ctx, "tcp", ":10224")
	if err != nil {
		t.Error(err)
		return
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	go func() {
		<-ctx.Done()
		grpcServer.Stop()
	}()

	for _, handler := range handlers {
		handler(grpcServer)
	}

	grpcServer.Serve(lis)
}

func dialWithHandler(t *testing.T, handler func(net1.Conn), grpcHandlers ...func(*grpc.Server)) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	wait := make(chan bool)
	go func(t *testing.T) {
		startGrpcServer(ctx, t, grpcHandlers...)
		close(wait)
	}(t)

	time.Sleep(time.Second)

	conn, err := Dialer.Dial(ctx, "127.0.0.1:10224")
	assert.Assert(t, err == nil)

	handler(conn)

	cancel()
	<-wait
}

func TestGrpcDailTimeOut(t *testing.T) {

	_, err := Dialer.Dial(context.Background(), "ddd")
	assert.Equal(t, err, context.DeadlineExceeded, "连接不可连接点，将返回超时")
}

func TestGrpcDialSuccess(t *testing.T) {

	dialWithHandler(t, func(c net1.Conn) {
	})
}

func TestPipeAutoLinkClient(t *testing.T) {

	dialWithHandler(t, func(c net1.Conn) {
		pp := c.(*Pipe)
		assert.Assert(t, pp.autoLinkClient() != nil)
		assert.Equal(t, pp.autoLinkClient(), pp.autoLinkClient(), "多次获取的client必须是同一个对象")
	})
}

func TestPipeRecvClosed(t *testing.T) {

	pp := newPipe(context.Background(), "test", nil)
	go func() {
		pp.Close()
	}()

	msg, err := pp.Recv()
	assert.Assert(t, msg == nil, "关闭后，收到的消息为nil")
	assert.Equal(t, err, io.EOF, "关闭后，返回的错误为io.EOF")
}

type testPipeServer struct {
	tg.UnimplementedPipeServer
}

func (tp *testPipeServer) Link(pls tg.Pipe_LinkServer) error {

	for {
		msg, err := pls.Recv()
		switch err {
		case io.EOF:
			return err
		case nil:
			name := string(msg.GetValue().MessageName().Name()) //ptypes.AnyMessageName(msg.GetValue())
			if name[:2] == "S_" {
				msg.Key += "_STREAM"
				pls.Send(msg)
			}
		}
	}
}

func (tp *testPipeServer) Donce(ctx context.Context, msg *pb.Message) (*pb.Message, error) {

	name := string(msg.GetValue().MessageName().Name()) //ptypes.AnyMessageName(msg.GetValue())
	name = name[:2]
	switch name {
	case "P_":
		msg.Key += "_POST"
	case "G_":
		msg.Key += "_GET"
	}
	return msg, nil
}

func TestPipeSendError(t *testing.T) {

	wait := make(chan bool)
	dialWithHandler(t, func(c net1.Conn) {

		msg := &pb.Message{Key: "Post"}
		err := c.Send(msg)
		assert.Assert(t, err != nil, "Message's Value is nil, should get error: message is nil")

		msg.Value, _ = anypb.New(&td.TEST_Example{})
		err = c.Send(msg)
		assert.Equal(t, err.Error(), "cann't operate message type undefined", "只能处理 G_ | P_ | S_ 开头的消息")

		close(wait)

	}, func(grpcServer *grpc.Server) {
		tg.RegisterPipeServer(grpcServer, &testPipeServer{})
	})
	<-wait
}

func TestPipeSend(t *testing.T) {

	lock := make(chan bool)
	// wait := make(chan bool)
	dialWithHandler(t, func(c net1.Conn) {
		var wait sync.WaitGroup
		wait.Add(3)

		t.Run("Pipe: Stream", func(t *testing.T) {
			msg := &pb.Message{Key: "Stream"}
			msg.Value, _ = anypb.New(&td.S_Example{})
			err := c.Send(msg)
			assert.Assert(t, err == nil, "")

			hostname := os.Getenv("HOSTNAME")

			recvMsg, err := c.Recv()
			assert.Assert(t, err == nil)
			assert.Equal(t, recvMsg.GetKey(), hostname+strconv.FormatUint(msg.GetSenderUid(), 10)+"_STREAM")

			wait.Done()
		})

		t.Run("Pipe: Get", func(t *testing.T) {
			msg := &pb.Message{Key: "Get"}
			msg.Value, _ = anypb.New(&td.G_Example{})
			err := c.Send(msg)
			assert.Assert(t, err == nil, "")

			msgRecv, err := c.Recv()
			assert.Assert(t, err == nil)
			assert.Equal(t, msgRecv.GetKey(), msg.GetKey()+"_GET")

			wait.Done()
		})

		t.Run("Pipe: Post", func(t *testing.T) {

			msg := &pb.Message{Key: "Post"}
			msg.Value, _ = anypb.New(&td.P_Example{})
			err := c.Send(msg)
			assert.Assert(t, err == nil, "")

			msgChan := make(chan *pb.Message)
			go func(t *testing.T) {
				msgRecv, err := c.Recv()
				assert.Assert(t, err == nil)
				msgChan <- msgRecv
			}(t)

			select {
			case <-msgChan:
				assert.Assert(t, false, "Post type message cann't receive callback message")
			case <-time.After(time.Second):
			}
			wait.Done()
		})

		wait.Wait()
		close(lock)

	}, func(grpcServer *grpc.Server) {
		tg.RegisterPipeServer(grpcServer, &testPipeServer{})
	})
	// wait.Wait()
	<-lock
}
