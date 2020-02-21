package token

import (
	// "github.com/golang/glog"
	"errors"
	"sync"
	"context"
)

//Builder 用于构建Token
type Builder struct {

	ctx 	context.Context				//NOTE: Builder 的基础Context
	Cancel 	context.CancelFunc			//NOTE: 需要关闭时调用，关闭所有token

	hub 	map[string]*Token 			//NOTE: map[tokenStr]*Token
	rw		sync.RWMutex				//NOTE: 读写锁
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
	token, _ = newToken(context.WithValue(b.ctx, ContextKey("token"), tokenStr))		//NOTE: 此处一定不会触发error
	b.hub[tokenStr] = token
	return
}

//Get get *Token and if cann't find the *Token, then Build one and return
func (b *Builder) Get(tokenStr string) (*Token, error) {

	if tk, ok := b.Find(tokenStr); ok {
		return tk, nil				//NOTE: already owned tk, back it 
	}
	return b.Build(tokenStr)		//NOTE: 
}


//Delete delete *Token with string key
func (b *Builder) Delete(tokenStr string) {

	b.rw.Lock()
	defer b.rw.Unlock()

	delete(b.hub, tokenStr)
}


var builder *Builder
var once sync.Once

//GetBuilder return singleton tokenBuilder object
func GetBuilder() *Builder {

	once.Do(func () {

		builder = &Builder {
			hub : make(map[string]*Token),
		}
		builder.ctx, builder.Cancel = context.WithCancel(context.Background())
		go func ()  {			
			<- builder.ctx.Done()
			builder = nil
		}()			//NOTE: 如果builder Cancel 被调用，则清空builder
	})
	return builder
}

