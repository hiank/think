// handle message recv by ws with rpc

package think

import (
	"context"

	"github.com/golang/glog"

	"github.com/hiank/think/core"
	"github.com/hiank/think/core/k8s"
	"github.com/hiank/think/core/rpc"
)

//rpcHandler 处理rpc请求消息
type rpcHandler struct {
	*core.Pool
	ctx         context.Context
	recvHandler core.MessageHandler
}

func newRPCHandler(ctx context.Context, recvHandler core.MessageHandler) *rpcHandler {

	return &rpcHandler{
		Pool:        core.NewPool(ctx),
		ctx:         ctx,
		recvHandler: recvHandler,
	}
}

//Handle 实现core.MessageHandler
//将msg转发到指定client
func (rh *rpcHandler) Handle(msg core.Message) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			glog.Warning(err)
		}
	}()
	name := k8s.TryServerNameFromPBAny(msg.GetValue())

	return <-rh.AutoOne(name, func() (msgHub *core.MessageHub) {

		client := rpc.NewClient(rh.ctx, name)
		msgHub = core.NewMessageHub(rh.ctx, core.MessageHandlerTypeFunc(client.Send))
		go func() {
			if cc, err := client.Dial(k8s.TryServiceURL(rh.ctx, k8s.TypeKubIn, name+"service", "grpc")); err == nil {
				msgHub.DoActive()
				rh.Listen(client, rh.recvHandler)
				cc.Close()
			}
			rh.Del(name)
		}()
		return
	}).Push(msg)
}
