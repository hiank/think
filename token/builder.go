package token

import (
	"errors"
	"sync"
	"context"
)

type tokenBuilder struct {

	ctx 	context.Context				//NOTE: Builder 的基础Context
	Cancel 	context.CancelFunc

	hub 	map[string]*Token 			//NOTE: map[tokenStr]*Token
}


func (tb *tokenBuilder) get(tokenStr string) *Token {

	token, ok := tb.hub[tokenStr]
	if !ok {

		token, _ = newToken(context.WithValue(tb.ctx, ContextKey("token"), tokenStr))		//NOTE: 此处一定不会触发error
		tb.hub[tokenStr] = token
	}
	return token
}

func (tb *tokenBuilder) delete(tokenStr string) {

	delete(tb.hub, tokenStr)
}


var builder *tokenBuilder
var mtx sync.RWMutex

//InitBuilder 获取单例的tokenBuilder
func InitBuilder(ctx context.Context) {

	mtx.Lock()
	if builder == nil {

		builder = &tokenBuilder{
			hub : make(map[string]*Token),
		}
		builder.ctx, builder.Cancel = context.WithCancel(ctx)
	}
	mtx.Unlock()
}

//Get 根据字符token 找到Token，如果不存在，则新建一个
func Get(tokenStr string) (token *Token, err error) {

	mtx.RLock()
	defer mtx.RUnlock()

	if builder == nil {
		err = errors.New("package token error : without initialized builder. please call InitBulider(context.Context) first")
		return
	}

	select {
	case <-builder.ctx.Done():
		err = builder.ctx.Err()
		builder.Cancel()
		builder = nil
	default:
		token = builder.get(tokenStr)
	}
	return
}

//CloseBuilder 清除
func CloseBuilder() {

	mtx.Lock()
	if builder != nil {
		builder.Cancel()
		builder = nil
	}
	mtx.Unlock()
}
