package pb

import (
	"strings"

	"github.com/hiank/think/set/codes"
	"google.golang.org/protobuf/types/known/anypb"
)

//Message Type 用于甄别message 需要用那种方式调用
const (
	TypeUndefined = iota //NOTE: 未定义的类型，表明出错了
	TypeGET              //NOTE: get消息，需要一个返回
	TypePOST             //NOTE: post消息
	TypeSTREAM           //NOTE: 流消息
	TypeMQ               //NOTE: 消息队列
)

// GetServerType 获得服务类型
func GetServeType(anyMsg *anypb.Any) (t int, err error) {

	protoName := anyMsg.MessageName().Name()
	if !protoName.IsValid() {
		return TypeUndefined, codes.Error(codes.ErrorAnyMessageIsEmpty)
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

func GetServeName(anyMsg *anypb.Any) (svcName string, err error) {
	name := anyMsg.MessageName().Name()
	if !name.IsValid() {
		return "", codes.Error(codes.ErrorAnyMessageIsEmpty)
	}
	key := string(name)
	key = key[strings.IndexByte(key, '_')+1:]
	svcName = strings.ToLower(key)
	return
}
