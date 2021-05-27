package pool

import (
	"context"
	"io"
	"sync"

	"github.com/hiank/think/set/codes"
	"google.golang.org/protobuf/proto"
)

//Hub 存储资源
type Hub struct {
	ctx     context.Context
	handler Handler
	mux     sync.RWMutex //NOTE: 用于handler设置
	list    *ListMux
	limit   *LimitMux
	io.Closer
}

//NewHub 创建新的Hub
func NewHub(ctx context.Context, maxLoop int) *Hub {
	hub := &Hub{
		ctx:   ctx,
		list:  NewListMux(),
		limit: &LimitMux{max: maxLoop},
	}
	return hub
}

func (hub *Hub) asyncLoopWork(res proto.Message) {
	if !hub.limit.Retain() {
		hub.list.Push(res)
		return
	}

	go func(res proto.Message) {
		for res != nil {
			hub.handler.Handle(res)
			res = hub.list.Shift()
		}
		hub.limit.Release()
	}(res)
}

//Push 将资源送入Hub
func (hub *Hub) Push(res proto.Message) {
	if hub.curHandler() == nil || hub.limit.Locked() {
		hub.list.Push(res)
		return
	}
	hub.asyncLoopWork(res)
}

//curHandler 安全读取当前的handler
func (hub *Hub) curHandler() Handler {
	hub.mux.RLock()
	defer hub.mux.RUnlock()

	return hub.handler
}

//TrySetHandler set hanlder, if handler is nil, panic
func (hub *Hub) TrySetHandler(handler Handler) {
	if handler == nil {
		panic(codes.PanicNilHandler)
	}
	hub.mux.Lock()
	defer hub.mux.Unlock()

	switch {
	case hub.handler != nil: //NOTE: cannot reset handler
		panic(codes.PanicExistedHandler)
	case hub.limit.Locked(): //NOTE: limitMux's max should be number > 0, so when first TrySetHandler, limitMux cann't be locked
		panic(codes.PanicNonLimit)
	}

	hub.handler = handler
	for !hub.limit.Locked() {
		res := hub.list.Shift()
		if res == nil {
			break
		}
		hub.asyncLoopWork(res)
	}
}

//SetHandler 设置handler
func (hub *Hub) SetHandler(handler Handler) {
	hub.mux.Lock()
	defer hub.mux.Unlock()

	hub.handler = handler
	for !hub.limit.Locked() {
		res := hub.list.Shift()
		if res == nil {
			break
		}
		hub.asyncLoopWork(res)
	}
}

//Handler Hub的处理接口
type Handler interface {
	Handle(proto.Message) error
}

// //HandlerFunc 函数形式的Handler
// type HandlerFunc func(proto.Message) error

// //Handle 实现Handler的必要接口
// func (hf HandlerFunc) Handle(val proto.Message) error {
// 	return hf(val)
// }
