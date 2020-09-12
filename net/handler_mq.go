// handle message recv by ws with mq

package net

import (
	"sync"

	"github.com/hiank/think/utils/robust"
	"google.golang.org/protobuf/proto"

	"github.com/hiank/think/pb"

	"github.com/hiank/think/net/mq"
	"github.com/hiank/think/pool"
)

var _singleMQClient *mq.Client
var _singleMQClientOnce sync.Once

//TryMQ mq实例
//如果实例构建失败，会抛出异常，调用者注意处理异常
func TryMQ() *mq.Client {

	_singleMQClientOnce.Do(func() {
		_singleMQClient = mq.TryNewClient("")
	})
	return _singleMQClient
}

//mqHandle 处理消息队列
func mqHandle(msg *pool.Message) (err error) {

	defer robust.Recover(robust.Error, robust.ErrorHandle(func(e interface{}) {
		err = e.(error) //NOTE: TryMQ() 可能会抛出异常，此处针对这个可能性，设置返回值
	}))

	data, err := proto.Marshal(msg.GetData())
	robust.Panic(err)

	name, err := pb.AnyMessageNameTrimed(msg.GetData())
	robust.Panic(err)

	return TryMQ().Publish(name, data)
}
