package net

import (
	"errors"

	"github.com/hiank/think/net/rpc"
	"github.com/hiank/think/net/ws"
	"github.com/hiank/think/pb"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/token"
)

//ServeRPC 启动一个k8s服务，同一个进程只能有一个k8s服务
func ServeRPC(ip string, msgHandler rpc.MessageHandler) error {

	return rpc.ListenAndServe(token.BackgroundLife().Context, ip, msgHandler)
}

//ServeWS 启动一个ws服务
func ServeWS(ip string, msgHandler pool.MessageHandler) error {

	return ws.ListenAndServe(token.BackgroundLife().Context, ip, msgHandler)
}

//ServeWSDefault 启动一个默认消息处理的ws服务
//默认的MessageHandler 根据消息名起始标志调用mq 或rpc 转发消息
func ServeWSDefault(ip string) error {

	return ServeWS(ip, wsRecvHandler(1))
}

//wsRecvHandler 处理ws服务收到的消息
type wsRecvHandler int

//Handle implement pool.MessageHandler
func (wh wsRecvHandler) Handle(msg *pool.Message) error {

	t, err := pb.GetServerType(msg.GetData())
	if err != nil {
		return err
	}
	switch t {
	case pb.TypeGET:
		fallthrough
	case pb.TypePOST:
		fallthrough
	case pb.TypeSTREAM:
		return rpcHandle(msg)
	case pb.TypeMQ:
		return mqHandle(msg)
	}
	name, err := pb.AnyMessageNameTrimed(msg.GetData())
	if err == nil {
		err = errors.New("no method handle message: " + name)
	}
	return err
}
