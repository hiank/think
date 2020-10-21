package think

import (
	"context"
	"errors"

	"github.com/hiank/think/core"
	"github.com/hiank/think/core/mq"
	"github.com/hiank/think/core/pb"
	"github.com/hiank/think/core/rpc"
	"github.com/hiank/think/core/ws"
)

//ServeRPC 启动一个k8s服务，同一个进程只能有一个k8s服务
func ServeRPC(ip string, msgHandler rpc.ReadHandler) error {

	return rpc.ListenAndServe(context.Background(), ip, msgHandler)
}

//ServeWS 启用一个ws服务
func ServeWS(ip string, msgHandler core.MessageHandler) error {

	return ws.ListenAndServe(context.Background(), ip, msgHandler)
}

//ServeWSDefault 启动一个ws服务 默认方式
func ServeWSDefault(ip string) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	ctx := context.Background()
	mqHandler, rpcHandler := &MQHandler{mq.TryNewClient("")}, NewRPCHandler(ctx, new(ws.Writer))
	return ServeWS(ip, core.MessageHandlerTypeFunc(func(msg core.Message) error {
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
			return rpcHandler.Handle(msg)
		case pb.TypeMQ:
			return mqHandler.Handle(msg) //mqHandle(msg)
		}
		name, err := pb.AnyMessageNameTrimed(msg.GetValue())
		if err == nil {
			err = errors.New("no method handle message: " + name)
		}
		return err
	}))
}
