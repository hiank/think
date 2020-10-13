package rpc

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/hiank/think/core"

	"github.com/hiank/think/core/pb"
	tg "github.com/hiank/think/core/rpc/pb"
	"github.com/hiank/think/settings"
	"google.golang.org/grpc"
)

type linkConn struct {
	key string
	ls  tg.Pipe_LinkServer
}

func (lc *linkConn) GetKey() string {

	return lc.key
}

func (lc *linkConn) Send(msg core.Message) error {

	return nil
}

func (lc *linkConn) Recv() (core.Message, error) {

	return nil, nil
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

	defer core.Recover(core.Warning)

	msg, err := ls.Recv()
	core.Panic(err)

	return s.Listen(&linkConn{key: msg.GetKey(), ls: ls}, s.handler)
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
			err = s.handler.HandlePost(req)
		}
	}
	return
}

var _singleServer *Server

//Writer 服务端写消息对象
type Writer int

//Handle 实现pool.MessageHandler
func (w Writer) Handle(msg core.Message) error {

	defer core.Recover(core.Fatal)
	if _singleServer == nil {
		core.Panic(errors.New("k8s server not started, please start a k8s server first. (use 'ListenAndServe' function to do this.)"))
	}
	return <-_singleServer.Push(msg)
}

// ListenAndServe start a PipeServer
func ListenAndServe(ctx context.Context, ip string, msgHandler ReadHandler) (err error) {

	defer func() {
		core.Recover(core.Fatal)
		_singleServer = nil
	}()

	if _singleServer != nil {
		err = errors.New("k8s server existed, cann't start another one")
		core.Panic(err)
	}

	lis, err := new(net.ListenConfig).Listen(ctx, "tcp", core.WithPort(ip, settings.GetSys().K8sPort))
	core.Panic(err)

	ctx, cancel := context.WithCancel(ctx)
	_singleServer = newServer(ctx, msgHandler)
	defer cancel() //NOTE: 清理这个服务，按需执行

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	tg.RegisterPipeServer(grpcServer, _singleServer)
	go func() {
		<-ctx.Done()
		grpcServer.Stop()
	}()
	return grpcServer.Serve(lis)
}
