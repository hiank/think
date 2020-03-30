package k8s

import (
	"container/list"
	"context"
	"sync"
	"time"

	tg "github.com/hiank/think/net/k8s/protobuf"
	"github.com/hiank/think/pb"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/settings"
	"github.com/hiank/think/token"
	"google.golang.org/grpc"
)


type contextClientKey string


//Client k8s 客户端，每一个服务对应一个Client，连接池不关心
type Client struct {

	ctx 		context.Context

	ccMtx 		sync.RWMutex				//NOTE: 
	ccList 		*list.List					//NOTE: 保存ClientConn

	msgMtx 		sync.Mutex					//NOTE: 用于等待发送的消息 的存取
	waitingMsgs map[string]*list.List		//NOTE: 保存等待发送的消息 map[tokenStr]messageList

	*pool.Pool 					//NOTE: 每个token 对应一个pipe
}


//newClient 构建新的 Client，service 包含端口号
func newClient(ctx context.Context) *Client {

	c := &Client{
		ccList: 		list.New(),
		waitingMsgs: 	make(map[string]*list.List),
	}
	ctx = context.WithValue(ctx, pool.CtxKeyRecvHandler, ctx.Value(CtxKeyClientHubRecvHandler))
	ctx = context.WithValue(ctx, pool.CtxKeyConnBuilder, c)
	c.ctx, c.Pool = ctx, pool.NewPool(ctx)
	return c
}


//BuildAndSend for pool.ConnHandler
//when pool's ConnHub cann't find the *Conn, it would call this api in an new goroutine
func (c *Client) BuildAndSend(msg *pb.Message) {

	if c.pushWaitingMsg(msg) {
		return
	}
	
	cc := c.findCC(msg)
	if cc == nil {
		name, _ := pb.GetServerKey(msg.GetData())
		cc, _ = c.dial(name)
	}

	tok, _ := token.GetBuilder().Get(msg.GetToken())
	pipe := &Pipe{ctx: tok, pipe: tg.NewPipeClient(cc)}
	go c.sendWaitingMsgsWithPipe(msg.GetToken(), pipe)
	c.AddConn(pool.NewConn(tok, pipe))
}

//findCC find *grpc.ClientConn existed
func (c *Client) findCC(msg *pb.Message) (cc *grpc.ClientConn) {

	c.ccMtx.RLock()
	defer c.ccMtx.RUnlock()

	if c.ccList.Len() >= settings.GetSys().GrpcGo {
		cc = c.ccList.Front().Value.(*grpc.ClientConn)
	}
	return
}

//pushWaitingMsg push waiting message to 'waitingMsgs'
//return if the message is not the first one
func (c *Client) pushWaitingMsg(msg *pb.Message) (notFirst bool) {

	c.msgMtx.Lock()
	defer c.msgMtx.Unlock()

	msgList, notFirst := c.waitingMsgs[msg.GetToken()]
	if !notFirst {
		msgList = list.New()
		c.waitingMsgs[msg.GetToken()] = msgList
	}
	msgList.PushBack(msg)
	return
}

//sendWaitingMsgsWithPipe
//send all waiting messages with key in an new goutine [controled by the caller]
//when the goroutine running, the *Conn already added to pool
//so new messsage should not push to the message list
func (c *Client) sendWaitingMsgsWithPipe(key string, pipe *Pipe) {

	c.msgMtx.Lock()
	defer c.msgMtx.Unlock()

	msgList, ok := c.waitingMsgs[key]
	if !ok {
		return
	}
	for element := msgList.Front(); element != nil; element = element.Next() {
		pipe.Send(element.Value.(*pb.Message))
	}
	delete(c.waitingMsgs, key)		//NOTE: 处理完了将消息列表删掉
} 


func (c *Client) dial(name string) (*grpc.ClientConn, error) {

	addr, err := ServiceNameWithPort(TypeKubIn, name + "-service", "grpc")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(c.ctx, 10*time.Second)
	defer cancel()
	grpcCC, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	c.ccMtx.Lock()
	c.ccList.PushBack(grpcCC)
	c.ccMtx.Unlock()
	return grpcCC, nil
}
