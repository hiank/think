package pool

import (
	"container/list"
	"sync"

	"google.golang.org/protobuf/proto"
)

//LimitMux 限制器，用于限制
type LimitMux struct {
	max int          //NOTE: 最大值
	mux sync.RWMutex //NOTE:
	cur int
}

//Locked 限制器是否出发
func (lm *LimitMux) Locked() bool {
	lm.mux.RLock()
	defer lm.mux.RUnlock()

	return lm.cur >= lm.max
}

//Retain 索引值+1
func (lm *LimitMux) Retain() (ok bool) {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	if ok = lm.cur < lm.max; ok {
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
func (lm *ListMux) Push(val proto.Message) {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	lm.cache.PushBack(val)
}

//Shift 取出最前面的数据
func (lm *ListMux) Shift() (val proto.Message) {
	lm.mux.Lock()
	defer lm.mux.Unlock()

	if ok := lm.cache.Len() > 0; ok {
		val = lm.cache.Remove(lm.cache.Front()).(proto.Message)
	}
	return
}
