package pb

import (
	"github.com/hiank/think/set/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var DefaultHandler = &liteHandler{handlerHub: make(map[string]MessageHandler)}

//liteHandler 简易消息处理者
type liteHandler struct {
	handlerHub map[string]MessageHandler
}

//Handle pool.MessageHandler 传入的消息
func (dh *liteHandler) Handle(msg proto.Message) (err error) {
	name := string(msg.ProtoReflect().Descriptor().Name())
	if name == "Any" {
		if msg, err = msg.(*anypb.Any).UnmarshalNew(); err != nil {
			return
		}
		name = string(msg.ProtoReflect().Descriptor().Name())
	}

	if handler, ok := dh.handlerHub[name]; ok {
		err = handler.Handle(msg)
	} else {
		err = codes.ErrorNoMessageHandler
	}
	return
}

//Register 注册处理方法
//非线程安全，建议只在初始化时调用
func (dh *liteHandler) Register(emptyVal proto.Message, handler MessageHandler) error {
	name := string(emptyVal.ProtoReflect().Descriptor().Name())
	if _, ok := dh.handlerHub[name]; ok {
		return codes.ErrorExistedMessageHandler
	}

	dh.handlerHub[name] = handler
	return nil
}

//MessageHandler 处理proto.Message
//NOTE: 如果是*any.Any消息，会先执行UnmarshalAny 然后Handle返回的消息
type MessageHandler interface {
	Handle(proto.Message) error
	// UnmarshalAny(*any.Any) proto.Message
}
