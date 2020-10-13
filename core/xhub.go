package core

import (
	"context"
	"fmt"
	"sync"
)

//xHub 用于存储管理发送端MessageHub
type xHub struct {
	hub map[string]*MessageHub //NOTE: map[tokenString]*MessageHub
	mtx sync.RWMutex
}

//newXHub 构建xHub
func newXHub(ctx context.Context) *xHub {

	ch := &xHub{
		hub: make(map[string]*MessageHub),
	}
	return ch
}

func (ch *xHub) Add(key string, msgHub *MessageHub) (err error) {

	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	if _, ok := ch.hub[key]; !ok {
		ch.hub[key] = msgHub
	} else {
		err = fmt.Errorf("Conn tokened %v existed", key)
	}
	return
}

func (ch *xHub) Del(key string) {

	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	delete(ch.hub, key)
}

func (ch *xHub) Get(key string) (*MessageHub, bool) {

	ch.mtx.RLock()
	defer ch.mtx.RUnlock()

	msgHub, ok := ch.hub[key]
	return msgHub, ok
}

//AutoOne 先判断是否包含某MessageHub；包含 则不执行；不包含 则加写锁，并再次判断是否包含 判断是否执行
//这个方法将确保hub中有且只有一个MessageHub，并返回这个MessageHub
func (ch *xHub) AutoOne(key string, call func() *MessageHub) (msgHub *MessageHub) {

	var ok bool
	if msgHub, ok = ch.Get(key); ok {
		return
	}
	ch.mtx.Lock()
	defer ch.mtx.Unlock()
	if msgHub, ok = ch.hub[key]; !ok {
		msgHub = call()
		ch.hub[key] = msgHub
	}
	return
}
