package pool

import (
	"context"
)

//contextKey type of context key
type contextKey int

//CtxKeyRecvHandler 用于处理收到的消息，当连接池中的连接收到消息后，传入这个key对应的chan，有pool来处理
var CtxKeyRecvHandler = contextKey(0)

//CtxKeyConnBuilder used for ConnHandler in Context
//when ConnHub's Handle cann't find the *Conn
//if the CtxKeyConnBuilder's Value existed, call the handler's BuildAndSend
var CtxKeyConnBuilder = contextKey(1)


//Pool 连接池
type Pool struct {

	*ConnHub			//NOTE: 处理建立的连接
	*MessageHub			//NOTE: 处理转发的消息

	ctx 	context.Context
	Close 	context.CancelFunc
}


//NewPool 构建Pool
//ctx must contained 'CtxKeyRecvHandler'
func NewPool(ctx context.Context) *Pool {

	ctx, cancel := context.WithCancel(ctx)
	connHub := NewConnHub(ctx)
	p := &Pool {
		ConnHub			: connHub,
		MessageHub 		: NewMessageHub(connHub),
		ctx 			: ctx,
		Close 			: cancel,
	}
	return p
}
