package pool

import (
	"github.com/golang/glog"
	"time"
	"context"
)


//****************************************************************//

// DefaultTimer Timer 的一个默认实现
type DefaultTimer struct {

	unix 		int64 					//NOTE: 最近收到消息的时间戳，用于判断超时
	interval 	int64					//NOTE: 超时时间
}

// Update 更新状态
func (c *DefaultTimer) Update() {

	c.unix = time.Now().Unix()				//NOTE: 更新时间戳
}


// SetInterval 设置超时时间
func (c *DefaultTimer) SetInterval(interval int64) {

	glog.Infoln("interval : ", interval)
	c.interval = interval
}

// TimeOut 判断连接是否超时，第一次超时尝试建立通信
func (c *DefaultTimer) TimeOut() (out bool) {

	t := time.Now()
	interval := t.Unix() - c.unix

	out = (interval > c.interval)			//NOTE: 设定为10分钟，实际需要配表
	return
}


//****************************************************************//

//DefaultIdentifier 默认验证
type DefaultIdentifier struct {

	key 	string
	token 	string
}

//NewDefaultIdentifier 创建一个默认Identifier
func NewDefaultIdentifier(k string, t string) Identifier {

	return &DefaultIdentifier{k, t}
} 

//GetKey Identifier实现
func (di *DefaultIdentifier) GetKey() string {

	return di.key
}

//GetToken Identifier实现
func (di *DefaultIdentifier) GetToken() string {

	return di.token
}


//****************************************************************//

//Token 用于提供唯一信息
type Token struct {

	key 	string
	token 	string

	ctx 	context.Context
	Cancel 	context.CancelFunc
}

//NewToken 创建一个新的Token对象
func NewToken(ctx context.Context, key string, token string) *Token {

	t := &Token {
		key 	: key,
		token 	: token,
	}
	t.ctx, t.Cancel = context.WithCancel(ctx)
	return t
}

//GetKey 获得key
func (t *Token) GetKey() string {

	return t.key
}

//ToString 获得token 字符串
func (t *Token) ToString() string {

	return t.token
}
