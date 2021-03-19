package pb

import (
	"github.com/hiank/think/pool"
	"github.com/hiank/think/set/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var LiteHandler = &liteHandler{handlerHub: make(map[string]pool.Handler)}

//liteHandler 简易消息处理者
type liteHandler struct {
	DefaultHandler pool.Handler            //通用Handler，如果注册的Handler不匹配，尝试使用此Handler
	handlerHub     map[string]pool.Handler //MessageHandler
}

//Handle pool.MessageHandler 传入的消息
//step1:尝试从Hub中找到Handler，如果存在，直接调用并返回
//step2:如果是Any或pb.Message类型，解出包含的消息信息，再尝试一次，如果存在，直接调用并返回
//step3:如果存在UniversalHandler，调用之并返回
//step4:返回无法处理的错误码
func (dh *liteHandler) Handle(msg proto.Message) (err error) {

	name := string(msg.ProtoReflect().Descriptor().Name())
	if handler, ok := dh.handlerHub[name]; ok {
		return handler.Handle(msg)
	}

	switch name {
	case "Message":
		anyName := msg.(*Message).GetValue().MessageName().Name()
		if anyName.IsValid() {
			if handler, ok := dh.handlerHub[string(anyName)]; ok {
				return handler.Handle(msg)
			}
		}
	case "Any":
		anyName := msg.(*anypb.Any).MessageName().Name()
		if !anyName.IsValid() {
			return codes.ErrorAnyMessageIsEmpty
		}
		if handler, ok := dh.handlerHub[string(anyName)]; ok {
			if msg, err = msg.(*anypb.Any).UnmarshalNew(); err == nil {
				return handler.Handle(msg)
			}
		}
	}

	if dh.DefaultHandler != nil {
		return dh.DefaultHandler.Handle(msg)
	}
	return codes.ErrorNoMessageHandler
}

//Register 注册处理方法
//非线程安全，建议只在初始化时调用
func (dh *liteHandler) Register(emptyVal proto.Message, handler pool.Handler) error {
	name := string(emptyVal.ProtoReflect().Descriptor().Name())
	if _, ok := dh.handlerHub[name]; ok {
		return codes.ErrorExistedMessageHandler
	}

	dh.handlerHub[name] = handler
	return nil
}
