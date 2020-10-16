package rpc

import (
	"context"
	"errors"
	"sync"

	"github.com/hiank/think/core"
	"github.com/hiank/think/core/pb"
	tg "github.com/hiank/think/core/rpc/pb"
)

//Pipe ConnHandler for client conn
type Pipe struct {
	ctx      context.Context
	pipe     tg.PipeClient
	key      string
	recvChan chan core.Message
	stream   tg.Pipe_LinkClient //NOTE:
	linkOnce *sync.Once
	limit    chan byte
}

func newPipe(ctx context.Context, key string, pipe tg.PipeClient) *Pipe {

	return &Pipe{
		ctx:      ctx,
		pipe:     pipe,
		key:      key,
		recvChan: make(chan core.Message),
		linkOnce: new(sync.Once),
		limit:    make(chan byte, 1),
	}
}

//GetKey 获取Pipe的关键字，用于匹配消息
func (p *Pipe) GetKey() string {

	return p.key
}

//Send 向k8s服务端发送消息
//about stream: 实际上，只需要有一个steam 就可以了，这个是某个token对应的pipe，每个token只需要使用一个steam 足够了
func (p *Pipe) Send(msg core.Message) (err error) {

	defer core.Recover(core.Warning)

	t, err := pb.GetServerType(msg.GetValue())
	core.Panic(err)

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
	core.Panic(err)
	return
}

//Recv 将消息转换到此接口
func (p *Pipe) Recv() (msg core.Message, err error) {

	if msg, ok := <-p.recvChan; ok {
		return msg, nil
	}
	return nil, errors.New("k8s client read chan closed")
}

func (p *Pipe) sendByLink(msg *pb.Message) (err error) {

	defer core.Recover(core.Warning, func(r interface{}) {
		p.stream.CloseSend()
		p.stream = nil
		p.linkOnce = new(sync.Once)
		err = r.(error)
	})

	// reqMsg := &pb.Message{Key: p.key} //NOTE: 这里需要将当前服务名获取到，这个p.key 就是token，与msg.GetKey() 是相同的值
	// reqMsg := &pb.Message{Key: p.key + "/" + msg.GetKey()} //NOTE: 如果仅仅使用token 作为key 的话，rpc 服务将无法区分不同调用者(比如多个微服务访问某个提供rpc的公共服务)
	p.linkOnce.Do(func() {
		if p.stream, err = p.pipe.Link(p.ctx); err != nil {
			return
		}
		core.Panic(p.stream.Send(&pb.Message{Key: p.key})) //NOTE:
		go p.loopReadFromLink(p.stream)
	})

	p.limit <- 0 //NOTE: 某个时刻只能执行一个Send，事实上发送端也只有一个Recv
	core.Panic(p.stream.Send(&pb.Message{Key: p.key, Value: msg.GetValue()}))
	<-p.limit
	return
}

func (p *Pipe) loopReadFromLink(stream tg.Pipe_LinkClient) {

L:
	for {
		select {
		case <-p.ctx.Done():
			break L
		default:
			msg, err := stream.Recv()
			if err != nil {
				break L
			}
			p.recvChan <- &pb.Message{Key: p.key, Value: msg.GetValue()} //NOTE: 读到的message 的key 是包含当前服务唯一标志的，需要转换下
		}
	}
}
