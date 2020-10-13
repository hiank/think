package pb

import (
	"github.com/golang/protobuf/ptypes"

	"strings"

	"github.com/golang/protobuf/ptypes/any"
)

//AnyMessageNameTrimed 处理后的any.Any 消息名，去掉可能包含的包名
func AnyMessageNameTrimed(anyMsg *any.Any) (messageName string, err error) {

	if messageName, err = ptypes.AnyMessageName(anyMsg); err == nil {
		if dotIdx := strings.LastIndexByte(messageName, '.'); dotIdx != -1 {
			messageName = messageName[dotIdx+1:] //NOTE: 协议可能包含包名，此处截掉包名
		}
	}
	return
}

//Message Type 用于甄别message 需要用那种方式调用
const (
	TypeUndefined = iota //NOTE: 未定义的类型，表明出错了
	TypeGET              //NOTE: get消息，需要一个返回
	TypePOST             //NOTE: post消息
	TypeSTREAM           //NOTE: 流消息
	TypeMQ               //NOTE: 消息队列
)

// GetServerType 获得服务类型
func GetServerType(anyMsg *any.Any) (t int, err error) {

	messageName, err := AnyMessageNameTrimed(anyMsg)
	if err != nil {
		return TypeUndefined, err
	}

	switch messageName[0] {
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
