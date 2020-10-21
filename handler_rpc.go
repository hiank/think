// handle message recv by ws with rpc

package think

import (
	"context"

	"github.com/golang/glog"

	"github.com/hiank/think/core"
	"github.com/hiank/think/core/k8s"
	"github.com/hiank/think/core/rpc"
)

//RPCHandler 处理rpc请求消息
type RPCHandler struct {
	*core.Pool
	ctx         context.Context
	recvHandler core.MessageHandler
}

//NewRPCHandler new RPCHandler instance
func NewRPCHandler(ctx context.Context, recvHandler core.MessageHandler) *RPCHandler {

	return &RPCHandler{
		Pool:        core.NewPool(ctx),
		ctx:         ctx,
		recvHandler: recvHandler,
	}
}

//Handle 实现core.MessageHandler
//将msg转发到指定client
func (rh *RPCHandler) Handle(msg core.Message) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			glog.Warning(err)
		}
	}()

	name := k8s.TryServerNameFromPBAny(msg.GetValue())
	return <-rh.AutoOne(name, func() *core.MessageHub {

		client := rpc.NewClient(rh.ctx, name)
		msgHub := core.NewMessageHub(rh.ctx, core.MessageHandlerTypeFunc(client.Send))
		go func() {
			if _, err := client.Dial(k8s.TryServiceURL(rh.ctx, k8s.TypeKubIn, name+"service", "grpc")); err != nil {
				rh.Del(name)
				glog.Warning(err)
				return
			}
			rh.Listen(client, rh.recvHandler)
		}()
		return msgHub
	}).Push(msg)
}
