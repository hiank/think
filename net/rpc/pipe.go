package rpc

import (
	"context"
	"errors"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/pb"
	tg "github.com/hiank/think/net/rpc/pb"
	"github.com/hiank/think/pool"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/protobuf/proto"
)

var hostname = os.Getenv("hostname")

//Dialer grpc连接器
//10秒钟超时
var Dialer = dialerFunc(func(ctx context.Context, target string) (conn net.Conn, err error) {

	cc, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBalancerName(roundrobin.Name), grpc.WithTimeout(time.Second*10)) //NOTE: block 为阻塞直到ready，insecure 为不需要验证的
	if err == nil {
		conn = newPipe(ctx, target, tg.NewPipeClient(cc))
	}
	return
})

type dialerFunc func(ctx context.Context, target string) (conn net.Conn, err error)

func (df dialerFunc) Dial(ctx context.Context, target string) (conn net.Conn, err error) {
	return df(ctx, target)
}

//Pipe ConnHandler for client conn
type Pipe struct {
	ctx        context.Context
	cancel     context.CancelFunc
	pipe       tg.PipeClient
	key        string
	recvChan   chan *pb.Message
	linkClient *net.Client
	once       sync.Once
}

func newPipe(ctx context.Context, key string, pipe tg.PipeClient) *Pipe {

	p := &Pipe{
		pipe:     pipe,
		key:      key,
		recvChan: make(chan *pb.Message, 8),
	}
	p.ctx, p.cancel = context.WithCancel(ctx)
	return p
}

//Key 获取Pipe的关键字，用于匹配消息
func (p *Pipe) Key() string {

	return p.key
}

//Send 向k8s服务端发送消息
//about stream: 实际上，只需要有一个steam 就可以了，这个是某个token对应的pipe，每个token只需要使用一个steam 足够了
func (p *Pipe) Send(msg *pb.Message) (err error) {

	t, err := pb.GetServerType(msg.GetValue())
	if err != nil {
		return
	}

	pbMsg := &pb.Message{Key: msg.GetKey(), Value: msg.GetValue()}
	switch t {
	case pb.TypeGET:
		if msg, err = p.pipe.Donce(p.ctx, pbMsg); err == nil {
			p.recvChan <- msg //NOTE: TypeGET消息转送到Recv 接口
		}
	case pb.TypePOST:
		_, err = p.pipe.Donce(p.ctx, pbMsg)
	case pb.TypeSTREAM:
		p.autoLinkClient().Push(pbMsg)
	default:
		err = errors.New("cann't operate message type undefined")
	}
	return
}

func (p *Pipe) autoLinkClient() *net.Client {

	p.once.Do(func() {
		p.linkClient = net.NewClient(p.ctx, dialerFunc(func(ctx context.Context, target string) (c net.Conn, err error) {

			lc, err := p.pipe.Link(ctx)
			if err == nil {
				c = &Conn{Sender: p.buildLinkSender(lc), Reciver: lc, Closer: net.CloserFunc(lc.CloseSend)}
			}
			return
		}), pool.HandlerFunc(func(val proto.Message) error {
			p.recvChan <- val.(*pb.Message)
			return nil
		}))
	})
	return p.linkClient
}

func (p *Pipe) buildLinkSender(lc tg.Pipe_LinkClient) net.Sender {

	return net.SenderFunc(func(msg *pb.Message) error {
		msg.Key = hostname + strconv.FormatUint(msg.GetSenderUid(), 10) //NOTE: 此处已经唯一 服务名。host用于grpc确定来源
		return lc.Send(msg)
	})
}

//Recv 将消息转换到此接口
func (p *Pipe) Recv() (*pb.Message, error) {

	if msg, ok := <-p.recvChan; ok {
		return msg, nil
	}
	return nil, io.EOF
}

//Close 关闭
func (p *Pipe) Close() error {

	p.cancel()
	close(p.recvChan)
	return nil
}
