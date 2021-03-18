package k8s

import (
	"strings"

	"github.com/hiank/think/set/codes"
	"google.golang.org/protobuf/types/known/anypb"
)

//TryServerNameFromPBAny 从protobuf Any 消息获取服务名
func TryServerNameFromPBAny(msg *anypb.Any) string {

	name := msg.MessageName().Name() //pb.AnyMessageNameTrimed(msg)
	if !name.IsValid() {
		panic(codes.ErrorAnyMessageIsEmpty)
	}

	key := string(name)
	key = key[strings.IndexByte(key, '_')+1:]
	return strings.ToLower(key)
}
