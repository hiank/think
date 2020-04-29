package pool

import (
	"container/list"
	"context"

	"github.com/hiank/think/pb"
	tk "github.com/hiank/think/token"
)

//MessageHub 集中处理*pb.Message，顺序handle
type MessageHub struct {
	handler MessageHandler   //NOTE: 处理消息
	req     chan *messageReq //NOTE: 待处理的消息管道
	lock    chan bool
}

//messageReq 用于管道传递
type messageReq struct {
	*Message
	err chan<- error
}

//Message 包含token的Message，便于维护生命周期
type Message struct {
	*pb.Message
	*tk.Token
}

//NewMessage 构建Message
func NewMessage(msg *pb.Message, tok *tk.Token) *Message {

	return &Message{
		Message: msg,
		Token:   tok,
	}
}

//NewMessageHub 构建新的 MessageHub
func NewMessageHub(ctx context.Context, handler MessageHandler) *MessageHub {

	mh := &MessageHub{
		handler: handler,
		req:     make(chan *messageReq),
		lock:    make(chan bool),
	}
	go mh.loop(ctx)
	return mh
}

func (mh *MessageHub) loop(ctx context.Context) {

	locked, hub, waited, wait := true, list.New(), false, make(chan bool)
	handle := func() {

		if locked || waited || (hub.Len() == 0) { //NOTE: 锁住 或 等待中 或 无消息
			return
		}
		waited = true
		go func(req *messageReq) {

			req.err <- mh.handler.Handle(req.Message)
			wait <- false
		}(hub.Remove(hub.Front()).(*messageReq))
	}

L:
	for {
		select {
		case <-ctx.Done():
			break L
		case locked = <-mh.lock:
			handle()
		case waited = <-wait:
			handle()
		case msg := <-mh.req:
			hub.PushBack(msg)
			handle()
		}
	}
}

//LockReq 用于锁定或解锁 Handle
func (mh *MessageHub) LockReq() chan<- bool {

	return mh.lock
}

//PushWithBack 将消息推送到hub中
//返回一个chan 用于接收结果
//这个好像特别危险，因为如果没有读管道的操作，就会阻塞进去
func (mh *MessageHub) PushWithBack(msg *Message, errChan chan<- error) {

	mh.req <- &messageReq{Message: msg, err: errChan}
}

//Push 无需监听处理结果的推送方法
func (mh *MessageHub) Push(msg *Message) {

	mh.PushWithBack(msg, nil)
}

//MessageHandler Message处理接口
type MessageHandler interface {
	Handle(*Message) error //NOTE: 处理Message
}

//MessageHandlerTypeFunc 函数形式的MessageHandler
type MessageHandlerTypeFunc func(*Message) error

//Handle MessageHub
func (mhf MessageHandlerTypeFunc) Handle(msg *Message) error {

	return mhf(msg)
}

//MessageHandlerTypeChan chan形式的MessageHandler
type MessageHandlerTypeChan chan<- *Message

//Handle MessageHub
func (mhc MessageHandlerTypeChan) Handle(msg *Message) error {

	mhc <- msg
	return nil
}
