package pool

import (
	"context"
	"sync"
)

//HubPool 客户端，每种类型的协议使用唯一的客户端
type HubPool struct {
	ctx   context.Context
	cache map[string]*Hub //HubGetter
	mux   sync.RWMutex
}

//NewHubPool 新建一个HubPool
func NewHubPool(ctx context.Context) *HubPool {

	return &HubPool{
		ctx:   ctx,
		cache: make(map[string]*Hub),
	}
}

//AutoHub 自动获取Hub，有则返回，无则创建
//如果是新建了Hub，isNew将为true，需要外部设置Handler
func (hp *HubPool) AutoHub(key string) (hub *Hub, isNew bool) {

	if hub = hp.GetHub(key); hub != nil {
		return
	}

	hp.mux.Lock()
	defer hp.mux.Unlock()

	if hub = hp.cache[key]; hub == nil { //NOTE: 写锁中，也要判断是否已经写入
		hub = NewHub(hp.ctx, &LimitMux{Max: 1000})
		hp.cache[key] = hub
		isNew = true
	}
	return
}

//GetHub 获取指定的Hub
func (hp *HubPool) GetHub(key string) *Hub {

	hp.mux.RLock()
	defer hp.mux.RUnlock()

	return hp.cache[key]
}

//Remove 删除指定Hub
func (hp *HubPool) Remove(key string) {

	hp.mux.Lock()
	defer hp.mux.Unlock()

	if hub, ok := hp.cache[key]; ok {
		if hub.Closer != nil {
			hub.Closer.Close()
		}
		delete(hp.cache, key)
	}
}

//RemoveAll 移除所有hub
//这个方法主要是希望将所有缓存的Hub执行一次Close
func (hp *HubPool) RemoveAll() {

	hp.mux.Lock()
	defer hp.mux.Unlock()

	for _, hub := range hp.cache {
		if hub.Closer != nil {
			hub.Closer.Close()
		}
	}
	hp.cache = make(map[string]*Hub)
}
