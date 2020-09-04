// handle message recv by ws with mq

package net

import (
	"errors"
	"sync"

	"github.com/hiank/think/pb"

	"github.com/hiank/think/net/mq"
	"github.com/hiank/think/pool"
)

var _singleMQClient *mq.Client
var _singleMQClientOnce sync.Once

//MQHandle 处理消息队列
func MQHandle(msg *pool.Message) (err error) {

	_singleMQClientOnce.Do(func() {
		_singleMQClient, err = mq.TryNewClient("")
	})
	if _singleMQClient == nil {
		return errors.New("cann't connect to nats")
	}
	data, err := pb.AnyEncode(msg.GetData())
	if err != nil {
		return
	}
	name, err := pb.AnyMessageNameTrimed(msg.GetData())
	if err != nil {
		return
	}
	return _singleMQClient.Publish(name, data)
}
