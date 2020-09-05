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

//MQInstance mq实例
func MQInstance() *mq.Client {

	_singleMQClientOnce.Do(func() {
		_singleMQClient, _ = mq.TryNewClient("")
	})
	return _singleMQClient
}

//mqHandle 处理消息队列
func mqHandle(msg *pool.Message) (err error) {

	mqClient := MQInstance()
	if mqClient == nil {
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
	return mqClient.Publish(name, data)
}
