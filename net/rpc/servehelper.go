package rpc

import (
	"context"
	"net/http"

	"net"

	tnet "github.com/hiank/think/net"
	"github.com/hiank/think/net/pb"
	tg "github.com/hiank/think/net/rpc/pb"
	"github.com/hiank/think/set/codes"
	"google.golang.org/grpc"
)

//PipeServer gprc服务
type PipeServer struct {
	tg.UnimplementedPipeServer
	ctx      context.Context
	connChan chan<- tnet.Conn
	handler  DonceHandler
}

//Link operate 'stream' type message
func (ps *PipeServer) Link(ls tg.Pipe_LinkServer) (err error) {

	var msg *pb.Message
	if msg, err = ls.Recv(); err == nil {
		ctx, cancel := context.WithCancel(ps.ctx)
		ps.connChan <- &Conn{
			Sender:  ls,
			Reciver: ls,
			key:     msg.GetKey(),
			Closer: tnet.CloserFunc(func() error {
				cancel()
				return nil
			}),
		}
		<-ctx.Done()
	}
	return
}

//Donce respond TypeGET | TypePOST message
func (ps *PipeServer) Donce(ctx context.Context, req *pb.Message) (res *pb.Message, err error) {

	select {
	case <-ps.ctx.Done():
		err = http.ErrServerClosed
	case <-ctx.Done():
		err = http.ErrServerClosed
	default:
		t, _ := pb.GetServeType(req.GetValue()) //NOTE: 此接口收到的消息必然是 TypeGET or TypePOST
		switch t {
		case pb.TypeGET:
			res, err = ps.handler.HandleGet(req)
		case pb.TypePOST:
			res, err = new(pb.Message), ps.handler.HandlePost(req) //NOTE: 如果返回的res 为nil，会导致一些问题，所以返回一个空消息
		default:
			err = codes.ErrorNotSupportType
		}
	}
	return
}

func (ps *PipeServer) mustEmbedUnimplementedPipeServer() {}

//ServeHelper websocket连接核心
type ServeHelper struct {
	tnet.Accepter
	ctx        context.Context
	close      context.CancelFunc
	pipeServer *PipeServer
	connChan   chan tnet.Conn
	addr       string
}

//NewServeHelper 新建一个ServeHelper
func NewServeHelper(ctx context.Context, addr string, handler DonceHandler) *ServeHelper {

	ctx, close := context.WithCancel(ctx)
	ch := make(chan tnet.Conn, 8)
	return &ServeHelper{
		ctx:        ctx,
		close:      close,
		connChan:   ch,
		addr:       addr,
		Accepter:   tnet.ChanAccepter(ch),
		pipeServer: &PipeServer{ctx: ctx, connChan: ch, handler: handler},
	}
}

//ListenAndServe 启动服务
func (helper *ServeHelper) ListenAndServe() error {

	lis, err := new(net.ListenConfig).Listen(helper.ctx, "tcp", helper.addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	go func() {
		<-helper.ctx.Done()
		grpcServer.Stop()
	}()

	tg.RegisterPipeServer(grpcServer, helper.pipeServer)
	return grpcServer.Serve(lis)
}

//Close 关闭
func (helper *ServeHelper) Close() error {

	close(helper.connChan)
	helper.close()
	return nil
}
