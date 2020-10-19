package core

import (
	"context"
	"fmt"

	"github.com/golang/glog"
)

//Pool 池，集中处理Conn及消息
type Pool struct {
	*xHub //NOTE: 不做隐藏，外部可以安全的调用Add Del Get 方法处理对应的
	ctx   context.Context
}

//NewPool 构建一个Pool
func NewPool(ctx context.Context) *Pool {
	return &Pool{
		xHub: newXHub(ctx),
		ctx:  ctx,
	}
}

//Listen 监听Conn
//@param handler 收到消息处理者
//conn.Close would called after loop Recv
func (pool *Pool) Listen(conn Conn, handler MessageHandler) error {

	pool.AutoOne(conn.GetKey(), func() *MessageHub {
		return NewMessageHub(pool.ctx, MessageHandlerTypeFunc(conn.Send))
	}).activate() //NOTE: 确保存在对应的MessageHub，并激活

L:
	for {
		select {
		case <-pool.ctx.Done():
			break L
		default:
			msg, err := conn.Recv()
			if err != nil {
				break L
			}
			if err := handler.Handle(msg); err != nil {
				glog.Warningf("message %v handle error\n", msg)
			}
		}
	}
	pool.Del(conn.GetKey()) //NOTE: 接收端检测到连接出了问题，删除连接
	return conn.Close()     //NOTE: conn的关闭，放在读协程中处理
}

//Push 通过此Pool 推送消息
func (pool *Pool) Push(msg Message) <-chan error {

	if msgHub, ok := pool.Get(msg.GetKey()); ok {
		return msgHub.Push(msg)
	}
	rlt := make(chan error)
	go func(rlt chan<- error) {
		rlt <- fmt.Errorf("cann't find messagehub for %v", msg.GetKey())
	}(rlt)
	return rlt
}
