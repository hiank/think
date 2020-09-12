// handle message recv by ws with rpc

package net

import (
	"context"
	"sync"

	"github.com/hiank/think/net/k8s"
	"github.com/hiank/think/net/ws"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/token"
	"github.com/hiank/think/utils/robust"

	"github.com/hiank/think/net/rpc"
)

type rpcVal struct {
	*pool.Message
	res chan<- error
}

//rpcLoop loop to operate msg
func rpcLoop(ctx context.Context, sChan <-chan *rpcVal) {

	hub := make(map[string]*rpc.Client)
L:
	for {
		select {
		case <-ctx.Done():
			break L
		case val := <-sChan:
			val.res <- rpcPush(ctx, val, hub)
		}
	}
}

func rpcPush(ctx context.Context, val *rpcVal, hub map[string]*rpc.Client) (err error) {

	defer robust.Recover(robust.Warning, func(e interface{}) {
		err = e.(error)
	})

	name, err := val.ServerName()
	robust.Panic(err)

	client, ok := hub[name]
	if !ok {
		addr := k8s.TryServiceURL(ctx, k8s.TypeKubIn, name+"service", "grpc")
		client = rpc.NewClient(context.WithValue(ctx, pool.CtxKeyRecvHandler, new(ws.Writer)), addr)
		hub[name] = client
	}
	client.Push(val.Message)
	return nil
}

var _rpcMsgChan chan *rpcVal
var _rpcLoopOnce sync.Once

//RPCHandle 处理消息
func rpcHandle(msg *pool.Message) error {

	_rpcLoopOnce.Do(func() {
		_rpcMsgChan = make(chan *rpcVal)
		go rpcLoop(token.BackgroundLife().Context, _rpcMsgChan)
	})
	errChan := make(chan error)
	_rpcMsgChan <- &rpcVal{
		Message: msg,
		res:     errChan,
	}
	return <-errChan
}
