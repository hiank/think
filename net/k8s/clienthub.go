package k8s

import (
	"context"
	"sync"

	"github.com/golang/glog"
	"github.com/hiank/think/pb"
	"github.com/hiank/think/pool"
)

type contextKey int

//CtxKeyClientHubRecvHandler ClientHub收到消息处理
var CtxKeyClientHubRecvHandler = contextKey(0)


//ClientHub 管理所有client
type ClientHub struct {

	mtx 	sync.RWMutex					//NOTE: 读写锁，因为client 可能会被异步创建，所以需要
	ctx 	context.Context
	hub 	map[string]*Client				//NOTE: 用于保存client
}

func newClientHub(ctx context.Context) *ClientHub {

	return &ClientHub{
		ctx 	: ctx,
		hub 	: make(map[string]*Client),
	}
}

//Handle 实现 pool.MessageHandler，用于处理转发到k8s 的消息
func (ch *ClientHub) Handle(msg *pool.Message) error {

	name, err := pb.GetServerKey(msg.GetData())
	if err != nil {
		glog.Warningln("cann't operate msg : ", err)
		return err
	}

	client, ok := ch.findClient(name)//ch.hub[name]
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


var _singleClientHub *ClientHub
var _singleClientHubOnce sync.Once

//ActiveClientHub 激活单例ClientHub
//第一次调用时，必须带有context 参数
//context 需包含一个 pool.MessageHandler，用于处理ClientHub收到的消息
func ActiveClientHub(ctx context.Context) bool {

	_singleClientHubOnce.Do(func () {

		_singleClientHub = newClientHub(ctx)
	})
	return _singleClientHub != nil
}