package k8s

import (
	"context"
	"sync"

	"github.com/golang/glog"
	"github.com/hiank/think/pool"
)

//ClientHub 管理所有client
type ClientHub struct {
	mtx sync.RWMutex //NOTE: 读写锁，因为client 可能会被异步创建，所以需要
	ctx context.Context
	hub map[string]*Client //NOTE: 用于保存client
}

//NewClientHub 构建一个全新的ClientHub 用于处理消息，将消息转发到相关k8s 服务中
//ctx should contain value keyed pool.CtxKeyRecvHandler, to operate the message recv from conn
func NewClientHub(ctx context.Context) *ClientHub {

	return &ClientHub{
		ctx: ctx,
		hub: make(map[string]*Client),
	}
}

//Handle 实现 pool.MessageHandler，用于处理转发到k8s 的消息
func (ch *ClientHub) Handle(msg *pool.Message) error {

	name, err := msg.ServerName()
	if err != nil {
		glog.Warningln("cann't operate msg : ", err)
		return err
	}

	client, ok := ch.findClient(name) //ch.hub[name]
	if !ok {
		ch.mtx.Lock()
		defer ch.mtx.Unlock()
		client = newClient(ch.ctx, name)
		ch.hub[name] = client
	}
	client.Push(msg)
	return nil
}

func (ch *ClientHub) findClient(name string) (client *Client, ok bool) {

	ch.mtx.RLock()
	defer ch.mtx.RUnlock()

	client, ok = ch.hub[name]
	return
}
