package k8s

import (
	"context"
	"sync"

	"github.com/golang/glog"
	"github.com/hiank/think/pb"
)

type contextKey int

//CtxKeyClientHubRecvHandler ClientHub收到消息处理
var CtxKeyClientHubRecvHandler = contextKey(0)
//CtxKeyServiceName 保存服务名，构建Client 时用到
var CtxKeyServiceName = contextKey(1)


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
func (ch *ClientHub) Handle(msg *pb.Message) error {

	name, err := pb.GetServerKey(msg.GetData())
	if err != nil {
		glog.Warningln("cann't operate msg : ", err)
		return err
	}

	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	client, ok := ch.hub[name]
	if !ok {
		client = newClient(context.WithValue(ch.ctx, CtxKeyServiceName, name))
		ch.hub[name] = client
	}
	client.Post(msg)
	return nil
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