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

	ctx 		context.Context
	Close 		context.CancelFunc

	*pool.MessageHub 						//NOTE: 用于统一处理需要通过当前client 处理的消
	pool 		*pool.Pool 					//NOTE: 每个token 对应一个pipe
}


//newClient 构建新的 Client，service 包含端口号
func newClient(ctx context.Context, name string) *Client {

	ctx = context.WithValue(ctx, pool.CtxKeyRecvHandler, ctx.Value(CtxKeyClientHubRecvHandler))

	c, msgChan := new(Client), make(chan *pool.Message)
	c.ctx, c.Close = context.WithCancel(ctx)
	c.pool = pool.NewPool(ctx)
	c.MessageHub = pool.NewMessageHub(ctx, pool.MessageHandlerTypeChan(msgChan))

	go c.loop(name, msgChan)
	return c
}


func (c *Client) loop(name string, msgChan <-chan *pool.Message) {

	cc, err := c.dial(name)
	if err != nil {
		return
	}
	defer cc.Close()

	c.LockReq() <- false		//NOTE: 解锁消息hub
	hubMap := make(map[string]*pool.MessageHub)
	L: for {

		select {
		case <-c.ctx.Done(): break L
		case msg :=<-msgChan:		//NOTE: 这个消息时候 c.hub 发过来的
			tok, _ := token.GetBuilder().Get(msg.GetToken())
			hub, ok := hubMap[tok.ToString()]
			if !ok {
				hub = pool.NewMessageHub(c.ctx, pool.MessageHandlerTypeFunc(c.pool.Post))
				hubMap[tok.ToString()] = hub
				go c.listenPipe(tok, tg.NewPipeClient(cc), hub.LockReq())
			}
			hub.Push(msg)
		}
	}
}


// //Push 推送消息，此处在调用pool.Post 之前，需要先确保msg 在pool中
// func (c *Client) Push(msg *pool.Message) {

// 	c.hub.Push(msg)		//NOTE: c.hub 在dial完成之前，一直处于locked 状态，加入的消息都会缓存
// }


//listenPipe 监听pipe
func (c *Client) listenPipe(tok *token.Token, pipeClient tg.PipeClient, lockReq chan<- bool) {

	pipe, added := &Pipe{ctx: tok, pipe: pipeClient}, make(chan interface{})
	go c.pool.Listen(tok, pipe, added)

	lockReq <- !(<-added).(bool)
}


//dial 建立，需要检测返回cc 的状态
func (c *Client) dial(name string) (cc *grpc.ClientConn, err error) {

	defer robust.Recover(robust.Warning)

	addr, err := ServiceNameWithPort(TypeKubIn, name + "-service", "grpc")
	robust.Panic(err)

	cc, err = grpc.DialContext(c.ctx, addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBalancerName(roundrobin.Name))		//NOTE: block 为阻塞知道ready，insecure 为不需要验证的
	robust.Panic(err)
	return
}
