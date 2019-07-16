package pool

import (
	"github.com/golang/glog"
	"time"
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

// //DefaultIdentifier 默认验证
// type DefaultIdentifier struct {

// 	key 	string
// 	token 	string
// }

// //NewDefaultIdentifier 创建一个默认Identifier
// func NewDefaultIdentifier(k string, t string) Identifier {

// 	return &DefaultIdentifier{k, t}
// } 

// //GetKey Identifier实现
// func (di *DefaultIdentifier) GetKey() string {

// 	return di.key
// }

// //GetToken Identifier实现
// func (di *DefaultIdentifier) GetToken() string {

// 	return di.token
// }


//****************************************************************//
