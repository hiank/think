package core

import (
	"context"
	"errors"
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

//AutoListen 自动化Listen，会创建并添加MessageHub，并在监听结束时删除此MessageHub
func (pool *Pool) AutoListen(conn Conn, handler MessageHandler) (err error) {

	if err = pool.Add(conn.GetKey(), NewMessageHub(pool.ctx, MessageHandlerTypeFunc(conn.Send))); err != nil {
		return //NOTE: 如果是多线程自动创建监听，可能会出现多次设置的问题(比如rpc中，收到数据查找相应client，Get方法是读锁，存在多个线程同时没找到MessageHub的可能，然后都启用监听)
	}
	err = pool.Listen(conn, handler)
	pool.Del(conn.GetKey()) //NOTE: 接收端检测到连接出了问题，删除连接
	return
}

//Listen 监听Conn
//@param handler 收到消息处理者
func (pool *Pool) Listen(conn Conn, handler MessageHandler) (err error) {

L:
	for {
		select {
		case <-pool.ctx.Done():
			err = errors.New("Conn's token Done")
			break L
		default:
			var msg Message
			if msg, err = conn.Recv(); err != nil {
				break L
			}
			if err := handler.Handle(msg); err != nil {
				glog.Warningf("message %v handle error\n", msg)
			}
		}
	}
	return
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
