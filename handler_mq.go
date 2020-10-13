// handle message recv by ws with mq

package think

import (
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/hiank/think/core"
	"github.com/hiank/think/core/mq"
	"github.com/hiank/think/core/pb"
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
func mqHandle(msg core.Message) (err error) {

	data, err := proto.Marshal(msg.GetValue())
	if err != nil {
		return
	}

	name, err := pb.AnyMessageNameTrimed(msg.GetValue())
	if err != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return TryMQ().Publish(name, data)
}
