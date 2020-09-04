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

//ServeWS 启动一个ws服务，同一个进程只能有一个ws服务
//收到的消息交给k8s ClientHub 来处理
func ServeWS(ip string) error {

	return ws.ListenAndServe(token.BackgroundLife().Context, ip, wsRecvHandler(1))
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
		return RPCHandle(msg)
	case pb.TypeMQ:
		return MQHandle(msg)
	}
	name, err := pb.AnyMessageNameTrimed(msg.GetData())
	if err == nil {
		err = errors.New("no method handle message: " + name)
	}
	return err
}
