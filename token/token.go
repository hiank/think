package token

import (
	"errors"
	"context"
)

//ContextKey 用于
type ContextKey string


//Token 用于提供唯一信息
type Token struct {

	context.Context
	cancel 	context.CancelFunc

	derived bool 				//NOTE: 用于标识是否是派生Token
}

//newToken 创建一个新的Token对象
func newToken(ctx context.Context) (t *Token, err error) {

	if ctx.Value(ContextKey("token")) == nil {
		err = errors.New("no ContextKey('token') Value in param ctx")
		return
	}

	t = &Token {
		derived : false,
	}
	t.Context, t.cancel = context.WithCancel(ctx)
	return
}


//Derive 派生一个Token，用于低等级的Token绑定的生命周期维护
func (t *Token) Derive() *Token {

	derivedToken := &Token{
		derived : true,
	}
	derivedToken.Context, derivedToken.cancel = context.WithCancel(t)
	return derivedToken
}

//ToString 获得token 字符串
func (t *Token) ToString() string {

	return t.Value(ContextKey("token")).(string)
}

//Cancel 清理Token
func (t *Token) Cancel() {

	t.cancel()								//NOTE: 这个重复调用是没有关系的, 参见token_test.TestContextCancel
	if !t.derived && builder != nil {		//NOTE: 不是派生Token，同时token生成器没有被消除

		GetBuilder().Delete(t.Value(ContextKey("token")).(string))
	}
}