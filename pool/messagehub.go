package pool

import (
	"container/list"
	"context"
	"sync"

	"github.com/hiank/think/pb"
	"github.com/hiank/think/settings"
	tk "github.com/hiank/think/token"
)

//MessageHub 集中处理*pb.Message，顺序handle
type MessageHub struct {
	ctx        context.Context
	handler    MessageHandler   //NOTE: 处理消息
	req        chan *messageReq //NOTE: 待处理的消息管道
	hub        *list.List       //NOTE: 用于缓存待处理的请求
	activeOnce sync.Once        //NOTE: 只会执行一次激活操作
	activeReq  chan bool
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
		ctx:       ctx,
		handler:   handler,
		hub:       list.New(),
		req:       make(chan *messageReq, settings.GetSys().MessageHubReqLen),
		activeReq: make(chan bool, 1),
	}
	go mh.loop()
	return mh
}

func (mh *MessageHub) loop() {

L:
	for {
		select {
		case <-mh.ctx.Done():
			break L
		case <-mh.activeReq:
			go mh.loopHandle()
			break L //NOTE: 激活之后放弃此循环，执行loopHandle循环
		case req := <-mh.req:
			mh.hub.PushBack(req)
		}
	}
}

func (mh *MessageHub) remoteReq(loopReqs chan *messageReq) (req chan<- *messageReq) {
	if mh.hub.Len() > 0 {
		req = loopReqs
	}
	return
}

func (mh *MessageHub) remoteVal() (msgReq *messageReq) {
	if mh.hub.Len() > 0 {
		msgReq = mh.hub.Front().Value.(*messageReq)
	}
	return
}

func (mh *MessageHub) fork(num *int, loopReqs <-chan *messageReq, completeReq chan<- bool) {

	maxNum := 100
	if *num < maxNum && len(loopReqs) > 0 {
		go func() {
			mh.syncHandle(loopReqs)
			select {
			case <-mh.ctx.Done():
			default:
				completeReq <- true
			}
		}()
		*num++
	}
}

func (mh *MessageHub) loopHandle() {

	num, loopReqs, completeReq := 0, make(chan *messageReq, 10), make(chan bool, 10)
	defer close(loopReqs) //NOTE: 此处需要关闭接收channel，避免此MessageHub关闭后，还未结束的处理goroutine中向此channel发送请求导致堵塞
	// defer close(completeReq)
L:
	for {
		select {
		case <-mh.ctx.Done():
			break L
		case req := <-mh.req:
			if mh.hub.Len() > 0 || len(loopReqs) == cap(loopReqs) { //NOTE: 如果负载已满
				mh.hub.PushBack(req)
				break
			}
			loopReqs <- req
			mh.fork(&num, loopReqs, completeReq)
		case mh.remoteReq(loopReqs) <- mh.remoteVal():
			mh.hub.Remove(mh.hub.Front())
			mh.fork(&num, loopReqs, completeReq)
		case <-completeReq:
			num--
		}
	}
}

//syncHandle 处理消息请求
//这个方法很可能会耗时较多，调用者视情况确定是否需要另起goroutine
func (mh *MessageHub) syncHandle(loopReqs <-chan *messageReq) {

L:
	for {
		select {
		case msgReq, ok := <-loopReqs:
			if !ok {
				break L
			}
			if err := mh.handler.Handle(msgReq.Message); err != nil {
				msgReq.err <- err
			}
		default:
			break L
		}
	}
}

//DoActive 激活此messagehub
func (mh *MessageHub) DoActive() {

	mh.activeOnce.Do(func() {
		mh.activeReq <- true
	})
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
