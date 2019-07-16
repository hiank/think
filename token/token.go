package token

import (
	"errors"
	"context"
)

//ContextKey 用于
type ContextKey string


//Token 用于提供唯一信息
type Token struct {

	ctx 	context.Context
	cancel 	context.CancelFunc

	derived bool 				//NOTE: 用于标识是否是派生Token
}

//newToken 创建一个新的Token对象
func newToken(ctx context.Context) (t *Token, err error) {

	// tokenStr := ctx.Value(ContextKey("token")).(string)
	if ctx.Value(ContextKey("token")) == nil {
		err = errors.New("no ContextKey('token') Value in param ctx")
		return
	}

	t = &Token {
		derived : false,
	}
	t.ctx, t.cancel = context.WithCancel(ctx)
	return
}

//WithValue 生成一个绑定Key 的Context
func (t *Token) WithValue(key, value string) context.Context {

	return context.WithValue(t.ctx, ContextKey(key), value)
}

// //GetContext 获得Token 的Context，用来处理Done
// func (t *Token) GetContext() context.Context {

// 	return t.ctx
// }

//Done 用于监听ctx.Done()
func (t *Token) Done() <-chan struct{} {

	return t.ctx.Done()
}

//Err 用于获得ctx.Err()
func (t *Token) Err() error {

	return t.ctx.Err()
}


//Derive 派生一个Token，用于低等级的Token绑定的生命周期维护
func (t *Token) Derive() *Token {

	derivedToken := &Token{
		derived : true,
	}
	derivedToken.ctx, derivedToken.cancel = context.WithCancel(t.ctx)
	return derivedToken
}

//ToString 获得token 字符串
func (t *Token) ToString() string {

	return t.ctx.Value(ContextKey("token")).(string)
}

//Cancel 清理Token
func (t *Token) Cancel() {

	t.cancel()								//NOTE: 这个重复调用是没有关系的, 参见token_test.TestContextCancel
	if !t.derived && builder != nil {		//NOTE: 不是派生Token，同时token生成器没有被消除

		mtx.Lock()
		builder.delete(t.ctx.Value(ContextKey("token")).(string))
		mtx.Unlock()
	}
}