package pb

import (
	"github.com/hiank/think/set/codes"
	"google.golang.org/protobuf/types/known/anypb"
)

// //AnyMessageNameTrimed 处理后的any.Any 消息名，去掉可能包含的包名
// func AnyMessageNameTrimed(anyMsg *anypb.Any) (messageName string, err error) {

// 	name := anyMsg.MessageName().Name()
// 	if name.IsValid() {
// 		messageName = string(name)
// 	}
// 	// if messageName, err = string(anyMsg.MessageName()); err == nil {
// 	// 	if dotIdx := strings.LastIndexByte(messageName, '.'); dotIdx != -1 {
// 	// 		messageName = messageName[dotIdx+1:] //NOTE: 协议可能包含包名，此处截掉包名
// 	// 	}
// 	// }
// 	return
// }

//Message Type 用于甄别message 需要用那种方式调用
const (
	TypeUndefined = iota //NOTE: 未定义的类型，表明出错了
	TypeGET              //NOTE: get消息，需要一个返回
	TypePOST             //NOTE: post消息
	TypeSTREAM           //NOTE: 流消息
	TypeMQ               //NOTE: 消息队列
)

// GetServerType 获得服务类型
func GetServerType(anyMsg *anypb.Any) (t int, err error) {

	protoName := anyMsg.MessageName().Name()
	if !protoName.IsValid() {
		return TypeUndefined, codes.ErrorAnyMessageIsEmpty
	}

	name := string(protoName)
	switch name[0] {
	case 'G':
		t = TypeGET
	case 'P':
		t = TypePOST
	case 'S':
		t = TypeSTREAM
	case 'M':
		t = TypeMQ
	default:
		t = TypeUndefined
	}
	return
}
