package k8s

import (
	"context"

	tg "github.com/hiank/think/net/k8s/protobuf"
	"github.com/hiank/think/pb"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/token"
	"github.com/hiank/think/utils/robust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)


type pipeState struct {

	ready 	bool 
	name 	string
}

//Client k8s 客户端，每一个服务对应一个Client，连接池不关心
type Client struct {

	ctx 		context.Context

	send 		chan *pb.Message				//NOTE: 发送消息管道
	readyCC 	chan *grpc.ClientConn

	pool 		*pool.Pool 					//NOTE: 每个token 对应一个pipe
}


//newClient 构建新的 Client，service 包含端口号
func newClient(ctx context.Context, name string) *Client {

	c := &Client{
		send: 		make(chan *pb.Message),
		readyCC:	make(chan *grpc.ClientConn),
	}
	ctx = context.WithValue(ctx, pool.CtxKeyRecvHandler, ctx.Value(CtxKeyClientHubRecvHandler))
	c.ctx, c.pool = ctx, pool.NewPool(ctx)

	ready := make(chan bool)
	go c.loop(name, ready)
	<-ready		//NOTE: 等待loop协程开启

	go c.dial(name)
	return c
}


func (c *Client) loop(name string, ready chan bool) {

	ready <- true
	var (
		cc *grpc.ClientConn
		hubMap = make(map[string]*pool.MessageHub)
		pp = make(chan *pipeState)
	)
	L: for {

		select {
		case <-c.ctx.Done(): break L
		case cc =<-c.readyCC:
			if cc == nil {
				break L
			}
			defer cc.Close()
			for key := range hubMap {
				tok, _ := token.GetBuilder().Get(key)
				go c.listenPipe(tok, tg.NewPipeClient(cc), pp)
			}
		case state := <-pp:		//NOTE: 监听 listenPipe 完成状态
			if state.ready {
				hubMap[state.name].LockChan() <- false
			} else {
				delete(hubMap, state.name)
			}
		case msg :=<-c.send:
			tok, _ := token.GetBuilder().Get(msg.GetToken())
			hub, ok := hubMap[tok.ToString()]
			if !ok {
				hub = pool.NewMessageHub(c.ctx, c)
				hubMap[tok.ToString()] = hub
				if cc != nil {
					go c.listenPipe(tok, tg.NewPipeClient(cc), pp)
				}
			}
			hub.Push(msg)
			// hub.InChan() <- msg
		}
	}
}


//Post 发送消息，此处在调用pool.Post 之前，需要先确保msg 在pool中
func (c *Client) Post(msg *pb.Message) {

	c.send <- msg
}


//Handle 处理消息
func (c *Client) Handle(msg *pb.Message) error {

	c.pool.Post(msg)
	return nil
}


//listenPipe 监听pipe
func (c *Client) listenPipe(tok *token.Token, pipeClient tg.PipeClient, state chan *pipeState) {

	pipe, added := &Pipe{ctx: tok, pipe: pipeClient}, make(chan interface{})
	go c.pool.Listen(tok, pipe, added)

	ps := &pipeState{name: tok.ToString()}
	ps.ready = (<-added).(bool)
	state <- ps
}


//dial 建立，需要检测返回cc 的状态
func (c *Client) dial(name string) {

	defer robust.Recover(robust.Warning)

	addr, err := ServiceNameWithPort(TypeKubIn, name + "-service", "grpc")
	robust.Panic(err)

	cc, err := grpc.DialContext(c.ctx, addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBalancerName(roundrobin.Name))		//NOTE: block 为阻塞知道ready，insecure 为不需要验证的
	robust.Panic(err)

	c.readyCC <- cc
}

