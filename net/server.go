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
type ChanAccepter <-chan Conn

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
	close       context.CancelFunc //NOTE: 如果服务停止，关闭相关Context
}

//NewServer 创建服务
func NewServer(ctx context.Context, helper ServeHelper, handler pool.Handler) *Server {

	ctx, close := context.WithCancel(ctx)
	return &Server{
		ctx:         ctx,
		close:       close,
		helper:      helper,
		hubPool:     pool.NewHubPool(ctx),
		recvHandler: handler,
	}
}

//ListenAndServe 启动服务
func (srv *Server) ListenAndServe() error {

	go srv.loopAccept()
	err := srv.helper.ListenAndServe()
	srv.hubPool.RemoveAll() //NOTE: 调用所有Hub的Close，关闭连接
	srv.helper.Close()
	return err
}

func (srv *Server) loopAccept() {

	loopAccept(srv.ctx, srv.helper.(Accepter), srv.handleAccept)
}

//Send 发送消息，找到相应数据集，处理消息发送
func (srv *Server) Send(msg *pb.Message) error {

	if msg == nil {
		return codes.ErrorNilValue
	}

	hub := srv.hubPool.GetHub(msg.GetKey())
	if hub == nil {
		return codes.ErrorNotExisted
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
