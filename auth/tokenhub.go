package auth

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"
)

type listValue struct {
	*Token
	*time.Ticker
}

//TokenHub 用于维护Token
type TokenHub struct {
	context.Context
	cancel   context.CancelFunc
	mtx      sync.RWMutex
	cache    map[string]*list.Element //NOTE: map[tkStr]*Token
	sortList *list.List               //NOTE: 对cache中的Token有序存储，提高检索效率(主要用于超时清理)
	resetReq chan string              //NOTE: 需要重置的Token 的string值 管道
}

//NewTokenHub 构建一个TokenHub
func NewTokenHub(ctx context.Context) *TokenHub {
	hub := &TokenHub{
		cache:    make(map[string]*list.Element),
		sortList: list.New(),
		resetReq: make(chan string, 10),
	}
	hub.Context, hub.cancel = context.WithCancel(ctx)
	go hub.loop()
	return hub
}

func (hub *TokenHub) loop() {
	interval := time.Second
L:
	for {
		select {
		case <-hub.Done():
			break L
		case tkStr := <-hub.resetReq:
			hub.resetToken(tkStr)
		case <-time.After(interval): //NOTE: 每隔一段时间执行一次清理，对超时token 执行根token 的Invalidate
			for hub.invalidateFont(interval) {
			}
		}
	}
}

//Reset 重置指定token 的滴答
//这个方法并不会等待重置成功
func (hub *TokenHub) Reset(tkStr string) {

	hub.resetReq <- tkStr
}

//TryGet 尝试获得Token，如果无法获得，抛出异常
func (hub *TokenHub) TryGet(tkStr string) *Token {
	if tkStr == "" {
		panic(errors.New("无效token字串"))
	}
	if val, ok := hub.cachedVal(tkStr); ok {
		return val.Value.(*listValue).Derive()
	}
	return hub.buildAndCache(tkStr).Derive()
}

func (hub *TokenHub) cachedVal(tkStr string) (*list.Element, bool) {
	hub.mtx.RLock()
	defer hub.mtx.RUnlock()

	val, ok := hub.cache[tkStr]
	return val, ok
}

//buildAndCache 构建并缓存Token
func (hub *TokenHub) buildAndCache(tkStr string) *Token {

	hub.mtx.Lock()
	defer hub.mtx.Unlock()

	if val, ok := hub.cache[tkStr]; ok { //NOTE: 避免并发问题，需要判断是否存在tk 先
		return val.Value.(*listValue).Token
	}

	tk := newToken(hub.Context, tkStr)
	hub.cache[tkStr] = hub.sortList.PushBack(&listValue{
		Token:  tk,
		Ticker: time.NewTicker(5 * time.Minute),
	})
	return tk
}

func (hub *TokenHub) resetToken(tkStr string) {

	hub.mtx.Lock()
	defer hub.mtx.Unlock()

	if val, ok := hub.cachedVal(tkStr); ok { //NOTE: resetToken 只会在loop中调用，并且删除操作也只会在loop中调用，所以不会出现安全问题

		ticker := val.Value.(*listValue).Ticker
		select {
		case <-ticker.C:
		default:
			ticker.Reset(5 * time.Minute)
		}
	}
}

func (hub *TokenHub) invalidateFont(interval time.Duration) (done bool) {

	hub.mtx.RLock()
	defer hub.mtx.Unlock()

	element := hub.sortList.Front()
	val := element.Value.(*listValue)
	select {
	case <-val.C: //NOTE: 此处表明，token 已经超时
		hub.sortList.Remove(element)
		delete(hub.cache, val.ToString())
		done = true
	default:
	}
	return
}
