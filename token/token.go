package token

import (
	"context"
	"errors"
	"time"

	"github.com/hiank/think/settings"
)

//contextKey 保存于Token的Context中的数据Key 类型
type contextKey int

//IdentityKey token 值 在Context中的key
var IdentityKey = contextKey(0)

//TimerKey Timer 在Context中的key
var TimerKey = contextKey(1)

//Token 用于提供唯一信息
type Token struct {

	context.Context
	cancel 	context.CancelFunc

	derived bool 					//NOTE: 是否是派生类型的token
}

//newToken 创建一个新的Token对象
func newToken(ctx context.Context) (t *Token, err error) {

	if ctx.Value(IdentityKey) == nil {
		err = errors.New("no 'IdentityKey' Value in param ctx")
		return
	}

	t = new(Token)
	t.Context, t.cancel = context.WithCancel(context.WithValue(ctx, TimerKey, time.NewTimer(time.Duration(settings.GetSys().TimeOut))))
	return
}


//Derive 派生一个Token，用于低等级的Token绑定的生命周期维护
func (t *Token) Derive() *Token {

	derivedToken := new(Token)
	derivedToken.Context, derivedToken.cancel = context.WithCancel(t)
	return derivedToken
}


//ResetTimer 重新设置定时器时间
func (t *Token) ResetTimer() {

	t.Value(TimerKey).(*time.Timer).Reset(time.Duration(settings.GetSys().TimeOut))
}

//ToString 获得token 字符串
func (t *Token) ToString() string {

	return t.Value(IdentityKey).(string)
}

//Cancel 清理Token
func (t *Token) Cancel() {

	t.cancel()								//NOTE: 这个重复调用是没有关系的, 参见token_test.TestContextCancel
	if !t.derived && _singleBuilder != nil {		//NOTE: 不是派生Token，同时token生成器没有被消除

		GetBuilder().Delete(t.Value(IdentityKey).(string))
	}
}