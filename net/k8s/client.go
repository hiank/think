package k8s

import (
	"context"

	tg "github.com/hiank/think/net/k8s/protobuf"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/token"
	"github.com/hiank/think/utils/robust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

//Client k8s 客户端，每一个服务对应一个Client，连接池不关心
type Client struct {
	*pool.MessageHub //NOTE: 用于统一处理需要通过当前client 处理的消
	ctx              context.Context
	Close            context.CancelFunc
	pool             *pool.Pool                  //NOTE: 每个token 对应一个pipe
	pipeHub          map[string]*pool.MessageHub //NOTE: 对每个token保存连接句柄
	rmReq            chan string                 //NOTE: 删除失效连接
}

//newClient 构建新的 Client，service 包含端口号
func newClient(ctx context.Context, name string) *Client {

	ctx, cancel := context.WithCancel(ctx)
	msgChan := make(chan *pool.Message)
	client := &Client{
		MessageHub: pool.NewMessageHub(ctx, messageHandlerTypeChan(msgChan)),
		ctx:        ctx,
		Close:      cancel,
		pool:       pool.NewPool(ctx),
		pipeHub:    make(map[string]*pool.MessageHub),
	}
	go client.loop(name, msgChan)
	return client
}

func (c *Client) loop(name string, msgChan <-chan *pool.Message) {

	cc, err := c.dial(name)
	if err != nil {
		return
	}
	defer cc.Close()
	c.DoActive() //NOTE: 解锁消息hub
L:
	for {
		select {
		case <-c.ctx.Done():
			break L
		case tokStr := <-c.rmReq:
			delete(c.pipeHub, tokStr)
		case msg := <-msgChan: //NOTE: 这个消息时候 c.hub 发过来的
			c.responseMessage(msg, cc)
		}
	}
}

//responseMessage 处理送进来的消息
func (c *Client) responseMessage(msg *pool.Message, cc *grpc.ClientConn) {

	tok, ok := token.GetBuilder().Find(msg.GetToken())
	if !ok { //NOTE: 找不到主Token，放弃处理
		return
	}
	hub, ok := c.pipeHub[tok.ToString()]
	if !ok {
		hub = pool.NewMessageHub(c.ctx, messageHandlerTypeFunc(c.pool.PostAndWait))
		c.pipeHub[tok.ToString()] = hub
		go func() {
			hub.DoActive() //NOTE: 这个hub带缓存，不会阻塞
			c.pool.Listen(tok.Derive(), &Pipe{ctx: tok, pipe: tg.NewPipeClient(cc)})
			c.rmReq <- tok.ToString() //NOTE: 连接结束后，删除
		}()
	}
	hub.Push(msg)
}

//dial 建立，需要检测返回cc 的状态
func (c *Client) dial(name string) (cc *grpc.ClientConn, err error) {

	defer robust.Recover(robust.Warning)

	addr, err := ServiceNameWithPort(c.ctx, TypeKubIn, name+"-service", "grpc")
	robust.Panic(err)

	cc, err = grpc.DialContext(c.ctx, addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBalancerName(roundrobin.Name)) //NOTE: block 为阻塞知道ready，insecure 为不需要验证的
	robust.Panic(err)
	return
}

type messageHandlerTypeChan chan<- *pool.Message

//Handle MessageHub
func (mhc messageHandlerTypeChan) Handle(msg *pool.Message) error {
	mhc <- msg
	return nil
}

//MessageHandlerTypeFunc 函数形式的MessageHandler
type messageHandlerTypeFunc func(*pool.Message) error

//Handle MessageHub
func (mhf messageHandlerTypeFunc) Handle(msg *pool.Message) error {
	return mhf(msg)
}
