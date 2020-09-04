package rpc

import (
	"context"
	"errors"

	"github.com/golang/glog"
	tg "github.com/hiank/think/net/rpc/protobuf"
	"github.com/hiank/think/pb"
)

//Pipe ConnHandler for client conn
type Pipe struct {
	ctx  context.Context
	pipe tg.PipeClient

	recvChan chan *pb.Message
	stream   tg.Pipe_LinkClient //NOTE:
}

func newPipe(ctx context.Context, pipe tg.PipeClient) *Pipe {

	return &Pipe{
		ctx:      ctx,
		pipe:     pipe,
		recvChan: make(chan *pb.Message),
	}
}

//Send 向k8s服务端发送消息
//about stream: 实际上，只需要有一个steam 就可以了，这个是某个token对应的pipe，每个token只需要使用一个steam 足够了
func (p *Pipe) Send(msg *pb.Message) (err error) {

	t, err := pb.GetServerType(msg.GetData())
	if err != nil {
		glog.Warningln(err)
		return
	}

	switch t {
	case pb.TypeGET:
		if msg, err = p.pipe.Donce(p.ctx, msg); err != nil {
			p.recvChan <- msg //NOTE: TypeGET消息转送到Recv 接口
		}
	case pb.TypePOST:
		_, err = p.pipe.Donce(p.ctx, msg)
	case pb.TypeSTREAM:
		err = p.sendByLink(msg)
	default:
		err = errors.New("cann't operate message type undefined")
	}
	if err != nil {
		glog.Warning("k8s client send message error : ", err)
	}
	return
}

//Recv 将消息转换到此接口
func (p *Pipe) Recv() (msg *pb.Message, err error) {

	var ok bool
	if msg, ok = <-p.recvChan; !ok {
		err = errors.New("k8s client read chan closed")
	}
	return
}

func (p *Pipe) sendByLink(msg *pb.Message) (err error) {

	if p.stream == nil {
		if p.stream, err = p.pipe.Link(p.ctx); err != nil {
			return
		}
		go p.loopReadFromLink(p.stream)
	}

	if err = p.stream.Send(msg); err != nil {
		p.stream.CloseSend()
		p.stream = nil
	}
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
				stream.CloseSend()
				break L
			}
			p.recvChan <- msg
		}
	}
}
