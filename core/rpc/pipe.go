package rpc

import (
	"context"
	"errors"
	"io"

	"github.com/hiank/think/core"
	"github.com/hiank/think/core/pb"
	tg "github.com/hiank/think/core/rpc/pb"
)

type pipeStream struct {
	key    string
	stream tg.Pipe_LinkClient
	limit  chan byte
}

func newPipeStream(key string) *pipeStream {

	return &pipeStream{
		key:   key,
		limit: make(chan byte, 1),
	}
}

func (ps *pipeStream) GetKey() string {

	return ps.key
}

func (ps *pipeStream) Send(msg core.Message) (err error) {

	ps.limit <- 0
	err = ps.stream.Send(&pb.Message{Key: msg.GetKey(), Value: msg.GetValue()})
	<-ps.limit
	return err
}

func (ps *pipeStream) Recv() (core.Message, error) {

	return ps.stream.Recv()
}

func (ps *pipeStream) Close() error {

	return ps.stream.CloseSend()
}

//Pipe ConnHandler for client conn
type Pipe struct {
	ctx      context.Context
	cancel   context.CancelFunc
	pipe     tg.PipeClient
	key      string
	recvChan chan core.Message
	linkPool *core.Pool
}

func newPipe(ctx context.Context, key string, pipe tg.PipeClient) *Pipe {

	ctx, cancel := context.WithCancel(ctx)
	return &Pipe{
		ctx:      ctx,
		cancel:   cancel,
		pipe:     pipe,
		key:      key,
		recvChan: make(chan core.Message),
		linkPool: core.NewPool(ctx),
	}
}

//GetKey 获取Pipe的关键字，用于匹配消息
func (p *Pipe) GetKey() string {

	return p.key
}

//Send 向k8s服务端发送消息
//about stream: 实际上，只需要有一个steam 就可以了，这个是某个token对应的pipe，每个token只需要使用一个steam 足够了
func (p *Pipe) Send(msg core.Message) (err error) {

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
		err = p.sendByLink(pbMsg)
	default:
		err = errors.New("cann't operate message type undefined")
	}
	return
}

//Recv 将消息转换到此接口
func (p *Pipe) Recv() (msg core.Message, err error) {

	if msg, ok := <-p.recvChan; ok {
		return msg, nil
	}
	return nil, io.EOF
}

//Close 关闭
func (p *Pipe) Close() error {

	p.cancel()
	return io.EOF
}

func (p *Pipe) sendByLink(msg *pb.Message) (err error) {

	return <-p.linkPool.AutoOne(p.key, func() *core.MessageHub {
		ps := newPipeStream(p.key)
		msgHub := core.NewMessageHub(p.ctx, core.MessageHandlerTypeFunc(ps.Send))
		go func() {
			if stream, err := p.pipe.Link(p.ctx); err == nil {
				if err = stream.Send(&pb.Message{Key: p.key}); err == nil {
					ps.stream = stream
					p.linkPool.Listen(ps, core.MessageHandlerTypeChan(p.recvChan))
					return
				}
				stream.CloseSend()
			}
			p.linkPool.Del(p.key) //NOTE: 如果连接出错，删除对应的MessageHub
		}()
		return msgHub
	}).Push(msg)
}
