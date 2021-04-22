package net

import (
	"context"
	"io"

	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/set/codes"
	"google.golang.org/protobuf/proto"
)

//ServeHelper 服务核心方法
type ServeHelper interface {
	ListenAndServe() error
	Accepter //NOTE: 务必保证Close后，Accept返回错误
	io.Closer
}

//Accepter 建立连接接口
type Accepter interface {
	Accept() (Conn, error)
}

//ChanAccepter chan方式的Accepter
type ChanAccepter chan Conn

//Accept 建立连接
func (ca ChanAccepter) Accept() (conn Conn, err error) {
	conn, ok := <-ca
	if !ok {
		err = io.EOF
	}
	return
}

//Server 服务
type Server struct {
	ctx         context.Context
	helper      ServeHelper
	hubPool     *pool.HubPool
	recvHandler pool.Handler       //NOTE: 处理收到的消息
	Close       context.CancelFunc //NOTE: 如果服务停止，关闭相关Context
}

//NewServer 创建服务
func NewServer(ctx context.Context, helper ServeHelper, handler pool.Handler) *Server {
	ctx, close := context.WithCancel(ctx)
	return &Server{
		ctx:         ctx,
		Close:       close,
		helper:      helper,
		hubPool:     pool.NewHubPool(ctx),
		recvHandler: handler,
	}
}

//ListenAndServe 启动服务
func (srv *Server) ListenAndServe() error {
	go loopAccept(srv.ctx, srv.helper.(Accepter), srv.handleAccept)
	go func(ctx context.Context) {
		<-ctx.Done()
		srv.helper.Close()
	}(srv.ctx)

	err := srv.helper.ListenAndServe()
	srv.hubPool.RemoveAll() //NOTE: 调用所有Hub的Close，关闭连接
	srv.Close()
	return err
}

//Send 发送消息，找到相应数据集，处理消息发送
func (srv *Server) Send(msg *pb.Message) error {

	select {
	case <-srv.ctx.Done():
		return codes.Error(codes.ErrorSrvClosed)
	default:
	}

	if msg == nil {
		return codes.Error(codes.ErrorNilValue)
	}

	hub := srv.hubPool.GetHub(msg.GetKey())
	if hub == nil {
		return codes.Error(codes.ErrorNotExisted)
	}
	hub.Push(msg)
	return nil
}

func (srv *Server) handleAccept(conn Conn) {

	hub, _ := srv.hubPool.AutoHub(conn.Key()) //NOTE: 这里暂时没考虑重复连接的问题，后续需要完善
	hub.SetHandler(pool.HandlerFunc(func(i proto.Message) error {
		return conn.Send(i.(*pb.Message))
	}))
	hub.Closer = conn.(io.Closer)

	go loopRecv(srv.ctx, conn.(Reciver), srv.recvHandler)
}

// type ConnHandler interface {
// 	Handle(Conn)
// }

//LiteServer 轻量级服务
//自定义 新连接及消息收发处理
type LiteServer struct {
	Context     context.Context
	ServeHelper ServeHelper
	ConnHandler func(Conn)
}

//ListenAndServe 启动服务
func (svc *LiteServer) ListenAndServe() error {
	go loopAccept(svc.Context, svc.ServeHelper.(Accepter), svc.ConnHandler)
	go func() {
		<-svc.Context.Done()
		svc.ServeHelper.Close()
	}()
	return svc.ServeHelper.ListenAndServe()
}
