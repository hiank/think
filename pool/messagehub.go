package pool

import (
	"container/list"
	"context"

	"github.com/hiank/think/pb"
)

//MessageHub 集中处理*pb.Message，顺序handle
type MessageHub struct {

	handler 	MessageHandler		//NOTE: 处理消息
	req 		chan *message		//NOTE: 待处理的消息管道
	lock 		chan bool
}

//message 用于管道传递
type message struct {

	*pb.Message
	err chan<- error
}


//NewMessageHub 构建新的 MessageHub
func NewMessageHub(ctx context.Context, handler MessageHandler) *MessageHub {

	mh := &MessageHub{
		handler: 	handler,
		req: 		make(chan *message),
		lock: 		make(chan bool),
	}
	go mh.loop(ctx)
	return mh
}


func (mh *MessageHub) loop(ctx context.Context) {

	locked, hub, waited, wait := true, list.New(), false, make(chan bool)
	handle := func ()  {

		if locked || waited || (hub.Len() == 0) {	//NOTE: 锁住 或 等待中 或 无消息
			return
		}
		waited = true
		go func (req *message) {

			req.err <- mh.handler.Handle(req.Message)
			wait <- false
		}(hub.Remove(hub.Front()).(*message))
	}

	L: for {
		select {
		case <-ctx.Done(): break L
		case locked = <-mh.lock: handle()
		case waited = <-wait: handle()
		case msg := <-mh.req:
			hub.PushBack(msg)
			handle()
		}
	}
}

//LockChan 用于锁定或解锁 Handle
func (mh *MessageHub) LockChan() chan<- bool {

	return mh.lock
}

//Push 将消息推送到hub中
//返回一个chan 用于接收结果
func (mh *MessageHub) Push(msg *pb.Message) <-chan error {

	err := make(chan error)
	mh.req <- &message{Message: msg, err: err}
	return err
}

// //InChan 加入消息管道
// func (mh *MessageHub) InChan() chan<- *pb.Message {

// 	return mh.req
// }


//MessageHandler Message处理接口
type MessageHandler interface {

	Handle(*pb.Message) error 		//NOTE: 处理Message
}