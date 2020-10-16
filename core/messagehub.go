package core

import (
	"container/list"
	"context"
	"sync"
	"time"
)

//messageReq 用于管道传递
type messageReq struct {
	Message
	rlt chan<- error
}

//MessageHub 集中处理Message，顺序handle
type MessageHub struct {
	ctx        context.Context
	handler    MessageHandler   //NOTE: 处理消息
	cacheReq   chan *messageReq //NOTE: 请求存储，对于无法即可处理的消息，需要缓存起来先，无缓存
	workChan   chan *messageReq //NOTE: 工作管道，每个工作线程尝试从这里读取请求，无缓存
	limit      chan byte        //NOTE: 限定工作线程，带缓冲，缓存值是工作线程的最大值
	activeOnce sync.Once        //NOTE: 只会执行一次激活操作
	notice     chan byte        //NOTE: 激活通知
}

//NewMessageHub 构建新的 MessageHub
func NewMessageHub(ctx context.Context, handler MessageHandler) *MessageHub {

	mh := &MessageHub{
		ctx:      ctx,
		handler:  handler,
		workChan: make(chan *messageReq),
		cacheReq: make(chan *messageReq),
		notice:   make(chan byte, 1),
	}
	go mh.loop()
	return mh
}

func (mh *MessageHub) loop() {

	cache := list.New()
L:
	for {
		workChan, val := mh.loopAutoWork(cache)
		select {
		case <-mh.ctx.Done():
			break L
		case workChan <- val: //NOTE: 尝试将缓存中第一个数据写入工作管道，成功的话，删掉第一个数据
			cache.Remove(cache.Front())
		case req := <-mh.cacheReq:
			cache.PushBack(req)
		case <-mh.notice:
			mh.limit = make(chan byte, 100)
		}
	}
}

//loopAutoWork 循环处理缓存中的消息，如果负载已满，则返回工作管道及缓存中的第一个数据
func (mh *MessageHub) loopAutoWork(cache *list.List) (chan<- *messageReq, *messageReq) {

	for cache.Len() > 0 {
		val := cache.Front().Value.(*messageReq)
		if !mh.autoWork(val) { //NOTE: 如果负载已满，则设置工作管道可写，并跳出此循环
			return mh.workChan, val
		}
		cache.Remove(cache.Front())
	}
	return nil, nil
}

func (mh *MessageHub) handle(req *messageReq) {

L:
	for {
		err := mh.handler.Handle(req.Message)
		ticker := time.NewTicker(time.Second)
		select {
		case req.rlt <- err:
		case <-ticker.C: //NOTE: 避免长时间阻塞
		}
		ticker.Reset(time.Second * 10)
		select {
		case req = <-mh.workChan:
		case <-ticker.C:
			break L
		}
	}
	<-mh.limit
}

//DoActive 激活此messagehub
func (mh *MessageHub) DoActive() {

	mh.activeOnce.Do(func() {
		mh.notice <- 0
	})
}

func (mh *MessageHub) autoWork(req *messageReq) (suc bool) {

	select {
	case mh.workChan <- req: //NOTE: 尝试写入工作管道
	default:
		select {
		case mh.limit <- 0: //NOTE: 尝试启动一个工作协程，如果工作协程数限制未到的话
			go mh.handle(req)
		default: //NOTE: 以上都无法完成，表明负载已满，暂时无法处理
			return false
		}
	}
	return true
}

//Push 推送消息，会返回一个错误chan，用于追踪推送结果
//如果不关心结果，则不处理返回的chan
//warning: 要及时处理返回的chan，超过1 * time.Millisecond 不监听 会导致chan丢失
func (mh *MessageHub) Push(msg Message) <-chan error {

	rlt := make(chan error)
	req := &messageReq{msg, rlt}
	if !mh.autoWork(req) { //NOTE: 如果负载已满，则将消息存入缓存
		mh.cacheReq <- req
	}
	return rlt
}

//MessageHandlerTypeChan chan type MessageHandler
type MessageHandlerTypeChan chan<- Message

//Handle MessageHub
func (mhc MessageHandlerTypeChan) Handle(msg Message) error {
	mhc <- msg
	return nil
}

//MessageHandlerTypeFunc 函数形式的MessageHandler
type MessageHandlerTypeFunc func(Message) error

//Handle MessageHub
func (mhf MessageHandlerTypeFunc) Handle(msg Message) error {
	return mhf(msg)
}
