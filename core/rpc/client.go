package rpc

import (
	"context"
	"fmt"

	"github.com/hiank/think/core"
	tg "github.com/hiank/think/core/rpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

//Client k8s 客户端，每一个服务对应一个Client，连接池不关心
type Client struct {
	context.Context
	Close    context.CancelFunc
	pipePool *core.Pool //NOTE: 每个Client 包含一组pipe
	key      string     //NOTE: 用于标识
	recv     chan core.Message
	cc       *grpc.ClientConn
}

//NewClient 构建新的 Client，service 包含端口号
//ctx should include pool.CtxKeyRecvHandler, that would used when client recv message
func NewClient(ctx context.Context, key string) *Client {

	ctx, cancel := context.WithCancel(ctx)
	client := &Client{
		Context:  ctx,
		Close:    cancel,
		pipePool: core.NewPool(ctx),
		key:      key,
		recv:     make(chan core.Message, 10),
	}
	return client
}

//Dial 建立，需要检测返回cc 的状态
func (c *Client) Dial(addr string) (cc *grpc.ClientConn, err error) {

	return grpc.DialContext(c.Context, addr, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name)) //NOTE: block 为阻塞知道ready，insecure 为不需要验证的
}

//GetKey 实现core.MessageHandler，返回的是服务名
func (c *Client) GetKey() string {

	return c.key
}

//Send 实现core.MessageHandler，对每个token 建立一个Pipe 用于发送消息
func (c *Client) Send(msg core.Message) error {

	key := msg.GetKey()
	return <-c.pipePool.AutoOne(key, func() (msgHub *core.MessageHub) {

		msgHub = core.NewMessageHub(c.Context, core.MessageHandlerTypeChan(c.recv))
		go func() {
			msgHub.DoActive()
			// c.pipePool.Listen(&Pipe{pipe: tg.NewPipeClient(c.cc)}, nil)
			c.pipePool.Listen(newPipe(c.Context, tg.NewPipeClient(c.cc)), nil)
		}()
		return
	}).Push(msg)
}

//Recv 实现core.MessageHandler，每个pipe 收到的消息，通过这个方法返回
func (c *Client) Recv() (core.Message, error) {

	if msg, ok := <-c.recv; ok {
		return msg, nil
	}
	return nil, fmt.Errorf("rpc client connected to %s: Recv error", c.key)
}
