package token

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hiank/think/settings"
)

//Builder 用于构建Token
type Builder struct {
	ctx    context.Context    //NOTE: Builder 的基础Context
	Cancel context.CancelFunc //NOTE: 需要关闭时调用，关闭所有token

	hub map[string]*Token //NOTE: map[tokenStr]*Token
	rw  sync.RWMutex      //NOTE: 读写锁
}

func newBuilder() *Builder {

	builder := &Builder{
		hub: make(map[string]*Token),
	}
	builder.ctx, builder.Cancel = context.WithCancel(context.Background())

	go builder.healthMonitoring()
	return builder
}

//healthMonitoring 监测状态
func (b *Builder) healthMonitoring() {

	interval := time.Duration(settings.GetSys().ClearInterval) * time.Second
L:
	for {
		select {
		case <-b.ctx.Done():
			_singleBuilder = nil
			break L
		case <-time.After(interval): //NOTE: 定时检查Token 超时
			for _, tok := range b.hub {
				select {
				case <-tok.Value(TimerKey).(*time.Timer).C: //NOTE: 如果某个tok 超时，关闭此token
					tok.Cancel()
				default:
				}
			}
		}
	}
}

//Find to get *Token with string key
func (b *Builder) Find(tokenStr string) (*Token, bool) {

	b.rw.RLock()
	defer b.rw.RUnlock()

	token, ok := b.hub[tokenStr]
	return token, ok
}

//Build build a *Token with string key
//Deprecated: Use Get instead
func (b *Builder) Build(tokenStr string) (token *Token, err error) {

	b.rw.Lock()
	defer b.rw.Unlock()

	if _, ok := b.hub[tokenStr]; ok {
		err = errors.New("token '" + tokenStr + "' existed in cluster")
		return
	}
	token, _ = newToken(context.WithValue(b.ctx, IdentityKey, tokenStr)) //NOTE: 此处一定不会触发error
	b.hub[tokenStr] = token
	return
}

//Get get *Token and if cann't find the *Token, then Build one and return
func (b *Builder) Get(tokenStr string) (*Token, error) {

	if tk, ok := b.Find(tokenStr); ok {
		return tk, nil //NOTE: already owned tk, back it
	}
	return b.Build(tokenStr) //NOTE:
}

//Delete delete *Token with string key
func (b *Builder) Delete(tokenStr string) {

	b.rw.Lock()
	defer b.rw.Unlock()

	delete(b.hub, tokenStr)
}

var _singleBuilder *Builder
var _singleBuilderOnce sync.Once

//GetBuilder return singleton tokenBuilder object
func GetBuilder() *Builder {

	_singleBuilderOnce.Do(func() {

		_singleBuilder = newBuilder()
	})
	return _singleBuilder
}
