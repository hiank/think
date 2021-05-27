package net

import (
	"context"

	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/pool"
	"k8s.io/klog/v2"
)

//Dialer 连接器，用于返回可写的对象
type Dialer interface {
	Dial(ctx context.Context, target string) (Conn, error)
}

//Client 客户端，每种类型的协议使用唯一的客户端
type Client struct {
	ctx         context.Context
	hubPool     *pool.HubPool
	dialer      Dialer
	recvHandler pool.Handler
}

//NewClient 新建一个Client
func NewClient(ctx context.Context, dialer Dialer, handler pool.Handler) *Client {

	return &Client{
		ctx:         ctx,
		dialer:      dialer,
		recvHandler: handler,
		hubPool:     pool.NewHubPool(ctx),
	}
}

//Push 使用Client，推送消息
//NOTE: 执行的步骤为
//1. 消息自动解析出服务名
//2. 从缓存中找到Hub
//3. 未找到Hub
//4. 新建一个Hub，并尝试建立连接
//5. 连接失败，删除Hub
func (client *Client) Push(msg *pb.Message) {

	hub := client.autoHub(msg.GetKey())
	hub.Push(msg)
}

func (client *Client) autoHub(key string) *pool.Hub {

	hub, needHandler := client.hubPool.AutoHub(key)
	if needHandler {
		go client.dialAndLoopRecv(hub, key)
	}
	return hub
}

func (client *Client) dialAndLoopRecv(hub *pool.Hub, key string) {

	conn, err := client.dialer.Dial(client.ctx, key)
	if err != nil {
		client.hubPool.Remove(key)
		klog.Warning("dial error: ", err)
		return
	}

	hub.SetHandler(HandlerFunc(func(msg *pb.Message) error {
		return conn.Send(msg)
	}))
	loopRecv(client.ctx, conn.(Reciver), client.recvHandler)
}
