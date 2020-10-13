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
	handler    MessageHandler //NOTE: 处理消息
	mtx        sync.RWMutex   //NOTE: 待处理消息缓存读写锁
	reqList    *list.List     //NOTE: 存储待处理消息
	notice     chan bool      //NOTE: 通知处理协程，有新消息
	activeOnce sync.Once      //NOTE: 只会执行一次激活操作
}

//NewMessageHub 构建新的 MessageHub
func NewMessageHub(ctx context.Context, handler MessageHandler) *MessageHub {

	mh := &MessageHub{
		ctx:     ctx,
		handler: handler,
		reqList: list.New(),
	}
	return mh
}

func (mh *MessageHub) loopHandle() {

	num, max, mtx := 0, 100, new(sync.Mutex)
L:
	for {
		select {
		case <-mh.ctx.Done():
			break L
		case <-mh.notice: //NOTE: notice缓存为1，调用者需要注意阻塞
			needNum := mh.safeLen() //NOTE: 这个地方使用读锁不会有问题，因为不存在多个goroutine 同时执行这个指令
			mtx.Lock()
			if needNum > max-num {
				needNum = max - num
			}
			num += needNum
			mtx.Unlock()
			for ; needNum > 0; needNum-- {
				go func() {
					mh.syncHandle()
					mtx.Lock()
					defer mtx.Unlock()
					num--
				}()
			}
		}
	}
}

func (mh *MessageHub) safePushBack(req *messageReq) {
	mh.mtx.Lock()
	defer mh.mtx.Unlock()
	mh.reqList.PushBack(req)
}

//safeShift 安全的方式删除并获取最前面的数据
//NOTE: 需要注意，如果是用读锁先获得list 长度，会有风险
//      当list长度为1，多个线程读取值，可能导致多个线程的到的数都为1，导致删除错误
func (mh *MessageHub) safeShift() (rlt *messageReq) {

	mh.mtx.Lock()
	defer mh.mtx.Unlock()
	if mh.reqList.Len() > 0 {
		rlt = mh.reqList.Remove(mh.reqList.Front()).(*messageReq)
	}
	return
}

func (mh *MessageHub) safeLen() int {
	mh.mtx.RLock()
	defer mh.mtx.RUnlock()
	return mh.reqList.Len()
}

//syncHandle 处理消息请求
//这个方法很可能会耗时较多，调用者视情况确定是否需要另起goroutine
func (mh *MessageHub) syncHandle() {

	req := mh.safeShift()
	for ; req != nil; req = mh.safeShift() {
		err := mh.handler.Handle(req.Message)
		select {
		case req.rlt <- err:
		case <-time.After(time.Millisecond): //NOTE: 避免长时间阻塞
		}
	}
}

//DoActive 激活此messagehub
func (mh *MessageHub) DoActive() {

	mh.activeOnce.Do(func() {
		mh.notice = make(chan bool, 1) //NOTE: 缓存为1，只需要接收一个通知
		mh.notice <- true              //NOTE: 通知处理消息，处理携程起来后，notice 一定是会被激活一次的。注意这个放在loopHandle进入的地方，可能会出现notice被Push写入，导致阻塞
		go mh.loopHandle()
	})
}

//Push 推送消息，会返回一个错误chan，用于追踪推送结果
//如果不关心结果，则不处理返回的chan
//warning: 要及时处理返回的chan，超过1 * time.Millisecond 不监听 会导致chan丢失
func (mh *MessageHub) Push(msg Message) <-chan error {

	rlt := make(chan error)
	mh.safePushBack(&messageReq{msg, rlt})
	select {
	case mh.notice <- true: //NOTE: 激活之前，mh.notice 为nil，会阻塞. 激活之后，只接收一个通知，其余忽略
	default: //NOTE: 激活之前，mh.notice阻塞，逻辑上会走这里
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
