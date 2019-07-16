package k8s

import (
	"net/http"
	"strconv"
	"bytes"
	"context"
	"net"
	"github.com/hiank/think/conf"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/pb"
	"google.golang.org/grpc"
	"github.com/golang/glog"
	tg "github.com/hiank/think/net/k8s/protobuf"
)


type server struct {

	MessageHandler
}

func (s *server) Link(ls tg.Pipe_LinkServer) (err error) {

	var msg *pb.Message
	if msg, err = ls.Recv(); err != nil {

		return err
	}
	// c := pool.NewDefaultConn(newConnHandler(pool.NewDefaultIdentifier(msg.GetKey(), msg.GetToken()), ls))
	c := pool.NewConn(msg.GetKey(), msg.GetToken(), newConnHandler(ls))
	// c.SetInterval(600)
	GetK8SPool().Push(c)
	return GetK8SPool().Listen(c)
}

func (s *server) Get(ctx context.Context, req *pb.Message) (res *pb.Message, err error) {

	glog.Infoln("do k8s get : ", req)
	select {
	case <-ctx.Done():
		err = http.ErrServerClosed
	default:
		res, err = s.HandleGet(req)
	}
	return
}

func (s *server) Post(ctx context.Context, msg *pb.Message) (rlt *tg.Void, err error) {

	select {
	case <-ctx.Done():
		err = http.ErrServerClosed
	default:
		err = s.HandlePost(msg)
	}
	return
}


// ListenAndServe start a PipeServer
func ListenAndServe(ctx context.Context, addr string, h MessageHandler) (err error) {

	k8sCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	InitK8SPool(k8sCtx, h.(pool.MessageHandler))

	var buffer bytes.Buffer
	buffer.WriteString(addr)
	buffer.WriteByte(':')
	buffer.WriteString(strconv.FormatInt(conf.GetSys().K8sPort, 10))

	addr = buffer.String()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	tg.RegisterPipeServer(grpcServer, &server{h})
	return grpcServer.Serve(lis)
}


//MessageHandler 服务器消息处理
type MessageHandler interface {

	pool.MessageHandler									//NOTE: 处理stream消息
	HandleGet(*pb.Message) (*pb.Message, error)		//NOTE: 处理Get消息
	HandlePost(*pb.Message) error 						//NOTE: 处理Post消息
}

//IgnoreGet 忽略Get 实现
type IgnoreGet int

//HandleGet 用于忽略HandleGet
func (i IgnoreGet) HandleGet(*pb.Message) (msg *pb.Message, err error) {

	return
}

//IgnorePost 忽略Post 实现
type IgnorePost int

//HandlePost 用于忽略HandlePost 方法
func (i IgnorePost) HandlePost(*pb.Message) (err error) {

	return
}

//IgnoreStream 忽略Stream 实现
type IgnoreStream int

//Handle 用于忽略Handle
func (i IgnoreStream) Handle(*pb.Message) (err error) {

	return
}
  