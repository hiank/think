package pool

import (
	"container/list"
	"sync"

	"github.com/hiank/think/pb"
	"github.com/hiank/think/settings"
)

//MessageHub message池，用于限制goroutine 数量来处理Message
type MessageHub struct {

	MessageHandler
	mtx 	*sync.Mutex			//NOTE: 处理message 增减
	queue 	*list.List 			//NOTE: message 队列
	gonum 	int					//NOTE: 当前goroutine 数量，用于限制goroutine 总数
}

//NewMessageHub 新建一个MessageHub
func NewMessageHub(handler MessageHandler) *MessageHub {

	return &MessageHub {
		MessageHandler 	: handler,
		mtx 			: new(sync.Mutex),
		queue 			: list.New(),
	}
}

//Post 新的待处理Message 传入其中 排队处理
func (mh *MessageHub) Post(msg *pb.Message) {

	mh.mtx.Lock()
	defer mh.mtx.Unlock()

	if mh.gonum < settings.GetSys().MessageGo {
		mh.gonum++
		go mh.do(msg)
		return
	}
	mh.queue.PushBack(msg)
}

func (mh *MessageHub) do(msg *pb.Message) {

	mh.Handle(msg)
	if msg = mh.shift(); msg != nil {
		mh.do(msg)
		return
	}
	mh.gonum--
}

func (mh *MessageHub) shift() (msg *pb.Message) {

	mh.mtx.Lock()
	defer mh.mtx.Unlock()

	if mh.queue.Len() > 0 {
		msg = mh.queue.Remove(mh.queue.Front()).(*pb.Message)
	}
	return
}


//MessageHandler Message处理接口
type MessageHandler interface {

	Handle(*pb.Message) error 		//NOTE: 处理Message
}