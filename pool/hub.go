package pool

import (
	"container/list"
	"context"
	"io"
	"sync"
)

//LimitMux 限制器，用于限制
type LimitMux struct {
	Max int          //NOTE: 最大值
	mux sync.RWMutex //NOTE:
	cur int
}

//Locked 限制器是否出发
func (lm *LimitMux) Locked() bool {
	lm.mux.RLock()
	defer lm.mux.RUnlock()

	return lm.cur >= lm.Max
}

//Retain 索引值+1
func (lm *LimitMux) Retain() (ok bool) {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	if ok = lm.cur < lm.Max; ok {
		lm.cur++
	}
	return
}

//Release 释放一个操作
func (lm *LimitMux) Release() {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	if lm.cur > 0 {
		lm.cur--
	}
}

//ListMux 线程安全list
type ListMux struct {
	cache *list.List
	mux   sync.Mutex
}

//NewListMux 创建一个ListMux
func NewListMux() *ListMux {

	return &ListMux{
		cache: list.New(),
	}
}

//Push 将数据安全的送到list中
func (lm *ListMux) Push(val interface{}) {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	lm.cache.PushBack(val)
}

//Shift 取出最前面的数据
func (lm *ListMux) Shift() (val interface{}) {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	if ok := lm.cache.Len() > 0; ok {
		val = lm.cache.Remove(lm.cache.Front())
	}
	return
}

//Hub 存储资源
type Hub struct {
	ctx     context.Context
	handler Handler
	mux     sync.RWMutex //NOTE: 用于handler设置
	limit   *LimitMux
	list    *ListMux
	io.Closer
}

//NewHub 创建新的Hub
func NewHub(ctx context.Context, limit *LimitMux) *Hub {

	return &Hub{
		ctx:   ctx,
		limit: limit,
		list:  NewListMux(),
	}
}

func (hub *Hub) asyncLoopWork(res interface{}) {

	if !hub.limit.Retain() {
		hub.list.Push(res)
		return
	}

	go func(res interface{}) {

		for res != nil {
			hub.handler.Handle(res)
			res = hub.list.Shift()
		}
		hub.limit.Release()
	}(res)
}

//loopWork 循环处理消息
func (hub *Hub) loopWork(res interface{}) {

	if !hub.limit.Retain() {
		hub.list.Push(res)
		return
	}

	for res != nil {
		hub.handler.Handle(res)
		res = hub.list.Shift()
	}
	hub.limit.Release()
}

//Push 将资源送入Hub
func (hub *Hub) Push(res interface{}) {

	if hub.curHandler() == nil || hub.limit.Locked() {
		hub.list.Push(res)
		return
	}

	hub.asyncLoopWork(res)
	// go hub.loopWork(res)
}

//curHandler 安全读取当前的handler
func (hub *Hub) curHandler() Handler {

	hub.mux.RLock()
	defer hub.mux.RUnlock()

	return hub.handler
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
	Handle(interface{}) error
}

//HandlerFunc 函数形式的Handler
type HandlerFunc func(interface{}) error

//Handle 实现Handler的必要接口
func (hf HandlerFunc) Handle(val interface{}) error {
	return hf(val)
}
