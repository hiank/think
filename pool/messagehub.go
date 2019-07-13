package pool

import (
	"github.com/golang/glog"
	"container/list"
	"sync"
	"github.com/hiank/think/pb"
)

//MessageHub message池，用于限制goroutine 数量来处理Message
type MessageHub struct {

	MessageHandler
	sync.RWMutex

	queue 	*list.List 			//NOTE: message 队列
	gonum 	int					//NOTE: 当前goroutine 数量
	MAXGO 	int					//NOTE: 最多goroutine 数量
}

//NewMessageHub 新建一个MessageHub
func NewMessageHub(handler MessageHandler, max int) *MessageHub {

	return &MessageHub {
		MessageHandler 	: handler,
		queue 			: list.New(),
		MAXGO 			: max,
	}
}

//Push 新的待处理Message 传入其中 排队处理
func (mh *MessageHub) Push(msg *pb.Message) {

	mh.Lock()
	defer mh.Unlock()

	glog.Infoln("messagehub MAXGO : ", mh.MAXGO)
	if mh.gonum < mh.MAXGO {

		mh.gonum++
		go mh.do(msg)
	} else {

		mh.queue.PushBack(msg)
	}	
}

func (mh *MessageHub) do(msg *pb.Message) {

	mh.Handle(msg)

	for front := mh.shift(); front != nil; front = mh.shift() {

		mh.Handle(front)
	}
	mh.gonum--
}

func (mh *MessageHub) shift() (msg *pb.Message) {

	mh.Lock()

	if len := mh.queue.Len(); len > 0 {

		front := mh.queue.Front()
		msg = mh.queue.Remove(front).(*pb.Message)
	}
	mh.Unlock()
	return
}


//MessageHandler Message处理接口
type MessageHandler interface {

	Handle(*pb.Message) error 		//NOTE: 处理Message
}