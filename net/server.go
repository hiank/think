package net

import (
	"context"
	"io"

	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/set/codes"
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

//liteServer provide 'pool' for Conn
type liteServer struct {
	ctx         context.Context
	hubPool     *pool.HubPool
	recvHandler pool.Handler
}

//handleConn add the conn at pool
func (svc *liteServer) handleConn(conn Conn) {
	hub, _ := svc.hubPool.AutoHub(conn.Key()) //NOTE: 这里暂时没考虑重复连接的问题，后续需要完善
	hub.SetHandler(HandlerFunc(func(msg *pb.Message) error {
		return conn.Send(msg)
	}))
	hub.Closer = conn //.(io.Closer)
	go loopRecv(svc.ctx, conn, svc.recvHandler)
}

//Send 发送消息，找到相应数据集，处理消息发送
func (svc *liteServer) Send(msg *pb.Message) (err error) {
	switch {
	case svc.ctx.Err() != nil:
		err = svc.ctx.Err()
	case msg == nil:
		err = codes.Error(codes.ErrorNilValue)
	default:
		if hub := svc.hubPool.GetHub(msg.GetKey()); hub != nil {
			hub.Push(msg)
		} else {
			err = codes.Error(codes.ErrorWorkerNotExisted)
		}
	}
	return err
}

//defaultLiteSender when use customize ConnHandler, the LiteSender will use this to return an error
var defaultLiteSender = SenderFunc(func(m *pb.Message) error { return codes.Error(codes.ErrorNonSupportLiteServe) })

//LiteSender when start a liteServer, someone can use LiteSender.Send(msg) for send the msg to client
var LiteSender Sender = defaultLiteSender

//ListenAndServe start listen serve
//NOTE: must pass one of WithConnHandler and WithRecvHandler, when pass WithConnHandler, WithRecvHandler will not work
func ListenAndServe(helper ServeHelper, opts ...ListenOption) error {
	if helper == nil {
		return codes.Error(codes.ErrorNonHelper)
	}

	dopts := newDefaultListenOptions()
	for _, opt := range opts {
		opt.apply(dopts)
	}
	ctx, cancel := context.WithCancel(dopts.ctx)
	defer cancel()

	var connHandler func(Conn)
	if dopts.connHandler == nil {
		if dopts.recvHandler == nil {
			return codes.Error(codes.ErrorNeedOneofConnRecvHandler)
		}
		svc := &liteServer{ctx: ctx, recvHandler: dopts.recvHandler, hubPool: pool.NewHubPool(ctx)}
		LiteSender = svc
		defer func() {
			svc.hubPool.RemoveAll()
			LiteSender = defaultLiteSender
		}()
		connHandler = svc.handleConn
	} else {
		connHandler = dopts.connHandler
	}

	go loopAccept(ctx, helper, connHandler)
	go func() {
		<-ctx.Done()
		helper.Close()
	}()
	return helper.ListenAndServe()
}
