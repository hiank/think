package rpc

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/golang/glog"

	"github.com/hiank/think/core"

	"github.com/hiank/think/core/pb"
	tg "github.com/hiank/think/core/rpc/pb"

	"github.com/hiank/think/settings"
	"google.golang.org/grpc"
)

type linkConn struct {
	key     string
	ls      tg.Pipe_LinkServer
	recvNum int
}

func (lc *linkConn) GetKey() string {

	return lc.key
}

func (lc *linkConn) Send(msg core.Message) error {

	return lc.ls.Send(&pb.Message{Key: msg.GetKey(), Value: msg.GetValue()})
}

func (lc *linkConn) Recv() (core.Message, error) {

	return lc.ls.Recv()
}

//Server k8s server
type Server struct {
	tg.UnimplementedPipeServer
	*core.Pool //NOTE: server 包含一个连接池，用于处理服务端的连接
	handler    ReadHandler
}

//newServer instantiate a Server
func newServer(ctx context.Context, msgHandler ReadHandler) *Server {

	return &Server{
		handler: msgHandler,
		Pool:    core.NewPool(ctx),
	}
}

//Link operate 'stream' type message
func (s *Server) Link(ls tg.Pipe_LinkServer) (err error) {

	msg, err := ls.Recv()
	if err != nil {
		return
	}
	return s.AutoListen(&linkConn{key: msg.GetKey(), ls: ls}, s.handler)
}

//Donce respond TypeGET | TypePOST message
func (s *Server) Donce(ctx context.Context, req *pb.Message) (res *pb.Message, err error) {

	select {
	case <-ctx.Done():
		err = http.ErrServerClosed
	default:
		t, _ := pb.GetServerType(req.GetValue()) //NOTE: 此接口收到的消息必然是 TypeGET or TypePOST
		switch t {
		case pb.TypeGET:
			res, err = s.handler.HandleGet(req)
		case pb.TypePOST:
			res, err = new(pb.Message), s.handler.HandlePost(req)
		}
	}
	return
}

var _singleServer *Server

//Writer 服务端写消息对象
type Writer int

//Handle 实现pool.MessageHandler
func (w Writer) Handle(msg core.Message) error {

	if _singleServer == nil {
		glog.Fatalln(errors.New("k8s server not started, please start a k8s server first. (use 'ListenAndServe' function to do this.)"))
	}
	return <-_singleServer.Push(msg)
}

// ListenAndServe start a PipeServer
func ListenAndServe(ctx context.Context, ip string, msgHandler ReadHandler) (err error) {

	if _singleServer != nil { //NOTE: 一个时刻只能存在一个服务
		return errors.New("k8s server existed, cann't start another one")
	}

	lis, err := new(net.ListenConfig).Listen(ctx, "tcp", core.WithPort(ip, settings.GetSys().K8sPort))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	_singleServer = newServer(ctx, msgHandler)
	defer func() {
		cancel() //NOTE: 清理这个服务，按需执行
		_singleServer = nil
	}()
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	tg.RegisterPipeServer(grpcServer, _singleServer)
	go func() {
		<-ctx.Done()
		grpcServer.Stop()
	}()
	return grpcServer.Serve(lis)
}
