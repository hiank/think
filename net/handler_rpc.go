// handle message recv by ws with rpc

package net

import (
	"context"
	"sync"

	"github.com/hiank/think/net/k8s"
	"github.com/hiank/think/net/ws"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/token"

	"github.com/hiank/think/net/rpc"
)

type rpcVal struct {
	*pool.Message
	res chan<- error
}

//rpcLoop loop to operate msg
func rpcLoop(ctx context.Context, sChan <-chan *rpcVal) {

	var val *rpcVal
	var err error
	hub := make(map[string]*rpc.Client)
L:
	for {
		select {
		case <-ctx.Done():
			break L
		case val = <-sChan:
			name, err := val.ServerName()
			if err != nil {
				break
			}
			client, ok := hub[name]
			if !ok {
				addr, err := k8s.ServiceNameWithPort(ctx, k8s.TypeKubIn, name+"service", "grpc")
				if err != nil {
					break
				}
				client = rpc.NewClient(context.WithValue(ctx, pool.CtxKeyRecvHandler, new(ws.Writer)), addr)
				hub[name] = client
			}
			client.Push(val.Message)
		}
		val.res <- err
	}
}

var _rpcMsgChan chan *rpcVal
var _rpcLoopOnce sync.Once

//RPCHandle 处理消息
func RPCHandle(msg *pool.Message) error {

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
