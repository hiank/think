package think

import (
	"context"
	"errors"

	"github.com/hiank/think/core"
	"github.com/hiank/think/core/pb"
	"github.com/hiank/think/core/rpc"
	"github.com/hiank/think/core/ws"
	"github.com/hiank/think/token"
	// "github.com/hiank/think/token"
)

//ServeRPC 启动一个k8s服务，同一个进程只能有一个k8s服务
func ServeRPC(ip string, msgHandler rpc.ReadHandler) error {

	return rpc.ListenAndServe(context.Background(), ip, msgHandler)
}

//ServeWS 启动一个ws服务
func ServeWS(ip string, msgHandler core.MessageHandler) error {

	return ws.ListenAndServe(token.BackgroundLife().Context, ip, msgHandler)
}

//ServeWSDefault 启动一个默认消息处理的ws服务
//默认的MessageHandler 根据消息名起始标志调用mq 或rpc 转发消息
func ServeWSDefault(ip string) error {

	return ServeWS(ip, core.MessageHandlerTypeFunc(HandleWS))
}

//HandleWS implement pool.MessageHandler
func HandleWS(msg core.Message) error {

	t, err := pb.GetServerType(msg.GetValue())
	if err != nil {
		return err
	}
	switch t {
	case pb.TypeGET:
		fallthrough
	case pb.TypePOST:
		fallthrough
	case pb.TypeSTREAM:
		// return rpcHandle(msg)
	case pb.TypeMQ:
		return mqHandle(msg)
	}
	name, err := pb.AnyMessageNameTrimed(msg.GetValue())
	if err == nil {
		err = errors.New("no method handle message: " + name)
	}
	return err
}
