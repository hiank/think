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
	cancel   context.CancelFunc
	resetReq chan bool //NOTE: 请求重置超时定时器
	derived  bool      //NOTE: 是否是派生类型的token
}

//newToken 创建一个新的Token对象
func newToken(ctx context.Context) (t *Token, err error) {

	if ctx.Value(IdentityKey) == nil {
		err = errors.New("no 'IdentityKey' Value in param ctx")
		return
	}
	ctx, cancel := context.WithCancel(ctx)
	t = &Token{
		Context:  ctx,
		cancel:   cancel,
		resetReq: make(chan bool, 10), //NOTE:带缓存，避免请求重置阻塞
	}
	go t.healthMonitoring()
	return
}

//healthMonitoring 健康监测
//处理超时，及context Done
func (t *Token) healthMonitoring() {

	duration := time.Millisecond * time.Duration(settings.GetSys().TimeOut)
	timer := time.NewTimer(duration) //NOTE:超时定时器
	defer timer.Stop()
L:
	for {
		select {
		case <-t.Done(): //NOTE: 关闭
			break L
		case <-timer.C: //NOTE: 超时
			t.Cancel()
			break L
		case <-t.resetReq:
			timer.Reset(duration)
		}
	}
}

//Derive 派生一个Token，用于低等级的Token绑定的生命周期维护
//此token可主动关闭，但不具备超时关闭的能力
func (t *Token) Derive() *Token {

	ctx, cancel := context.WithCancel(t)
	return &Token{
		Context: ctx,
		cancel:  cancel,
		derived: true,
	}
}

//ResetTimer 重新设置定时器时间
func (t *Token) ResetTimer() {

	t.resetReq <- true
}

//ToString 获得token 字符串
func (t *Token) ToString() string {

	return t.Value(IdentityKey).(string)
}

//Cancel 清理Token
func (t *Token) Cancel() {

	t.cancel()                             //NOTE: 这个重复调用是没有关系的, 参见token_test.TestContextCancel
	if !t.derived && GetBuilder() != nil { //NOTE: 不是派生Token，同时token生成器没有被消除
		GetBuilder().removeReq() <- t
	}
}
