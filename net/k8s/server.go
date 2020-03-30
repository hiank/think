package k8s

import (
	"context"
	"net"
	"net/http"

	"github.com/golang/glog"
	tg "github.com/hiank/think/net/k8s/protobuf"
	"github.com/hiank/think/pb"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/settings"
	"github.com/hiank/think/token"
	"github.com/hiank/think/utils"
	"google.golang.org/grpc"
)

//Server k8s server
type Server struct {

	handler 	MessageHandler
	*pool.Pool			//NOTE: server 包含一个连接池，用于处理服务端的连接
}

//newServer instantiate a Server
func newServer(ctx context.Context, msgHandler MessageHandler) *Server {

	return &Server {
		handler : msgHandler,
		Pool 	: pool.NewPool(context.WithValue(ctx, pool.CtxKeyRecvHandler, msgHandler)),
	}
}

//Link operate 'stream' type message
func (s *Server) Link(ls tg.Pipe_LinkServer) (err error) {

	var msg *pb.Message
	if msg, err = ls.Recv(); err != nil {
		return
	}

	tok, err := token.GetBuilder().Get(msg.GetToken())
	if err != nil {
		glog.Warning("cann't get token from token.GetBuilder(): ", err)
		return
	}
	return s.AddConn(pool.NewConn(tok, ls))
}

//Donce respond TypeGET | TypePOST message
func (s *Server) Donce(ctx context.Context, req *pb.Message) (res *pb.Message, err error) {

	glog.Infoln("do k8s Donce : ", req)
	select {
	case <-ctx.Done():
		err = http.ErrServerClosed
	default:
		t, _ := pb.GetServerType(req.GetData())		//NOTE: 此接口收到的消息必然是 TypeGET or TypePOST
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
func (w Writer) Handle(msg *pb.Message) error {

	if _singleServer != nil {
		glog.Fatalf("k8s server not started, please start a k8s server first. (use 'ListenAndServe' function to do this.)")
	}
	_singleServer.Pool.Post(msg)
	return nil
}


// ListenAndServe start a PipeServer
func ListenAndServe(ctx context.Context, ip string, msgHandler MessageHandler) (err error) {

	if _singleServer != nil {
		glog.Fatal("k8s server existed, cann't start another one.")
	}

	lis, err := net.Listen("tcp", utils.WithPort(ip, settings.GetSys().K8sPort))
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	_singleServer = newServer(ctx, msgHandler)
	defer _singleServer.Close()

	tg.RegisterPipeServer(grpcServer, _singleServer)
	return grpcServer.Serve(lis)
}
