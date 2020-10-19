package rpc

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/golang/glog"

	"github.com/hiank/think/core"
	"github.com/hiank/think/core/pb"
	tg "github.com/hiank/think/core/rpc/pb"
)

//Pipe ConnHandler for client conn
type Pipe struct {
	ctx      context.Context
	cancel   context.CancelFunc
	pipe     tg.PipeClient
	key      string
	recvChan chan core.Message
	stream   tg.Pipe_LinkClient //NOTE:
	linkOnce *sync.Once
	limit    chan byte
}

func newPipe(ctx context.Context, key string, pipe tg.PipeClient) *Pipe {

	ctx, cancel := context.WithCancel(ctx)
	return &Pipe{
		ctx:      ctx,
		cancel:   cancel,
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
	return nil, errors.New("k8s client read chan closed")
}

//Close 关闭
func (p *Pipe) Close() error {

	p.cancel()
	return io.EOF
}

func (p *Pipe) sendByLink(msg *pb.Message) (err error) {

	defer func() {
		if r := recover(); r != nil {
			p.linkOnce = new(sync.Once)
			err = r.(error)
			glog.Warning(err)
		}
	}()

	// reqMsg := &pb.Message{Key: p.key} //NOTE: 这里需要将当前服务名获取到，这个p.key 就是token，与msg.GetKey() 是相同的值
	// reqMsg := &pb.Message{Key: p.key + "/" + msg.GetKey()} //NOTE: 如果仅仅使用token 作为key 的话，rpc 服务将无法区分不同调用者(比如多个微服务访问某个提供rpc的公共服务)
	p.linkOnce.Do(func() {
		if p.stream, err = p.pipe.Link(p.ctx); err != nil {
			panic(err)
		}
		if err = p.stream.Send(&pb.Message{Key: p.key}); err != nil {
			p.stream.CloseSend()
			p.stream = nil
			panic(err)
		}
		go p.loopReadFromLink()
	})

	p.limit <- 0 //NOTE: 某个时刻只能执行一个Send，事实上发送端也只有一个Recv
	err = p.stream.Send(&pb.Message{Key: p.key, Value: msg.GetValue()})
	<-p.limit
	return
}

//loopReadFromLink 循环从建立的流中接收数据
//1. ctx关闭的话，主动调用流的CloseSend
//2. 收消息出错的话，如果是io.EOF则直接退出，否则也会执行流的CloseSend
// 退出的时候，会关闭recvChan，从而触发Recv 返回错误，使外部Listen 返回，触发相应的删除操作
func (p *Pipe) loopReadFromLink() {

	defer close(p.recvChan)
L:
	for {
		select {
		case <-p.ctx.Done():
			break L
		default:
			msg, err := p.stream.Recv()
			switch err {
			case nil:
				p.recvChan <- &pb.Message{Key: p.key, Value: msg.GetValue()} //NOTE: 读到的message 的key 是包含当前服务唯一标志的，需要转换下
			case io.EOF:
				return
			default:
				break L
			}
		}
	}
	var err error
	if err = p.stream.CloseSend(); err == nil {
		_, err = p.stream.Recv()
	}
	if err != io.EOF {
		glog.Warning(err)
	}
}
