package k8s

import (
	"strings"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/hiank/think/net/pb"
)

//TryServerNameFromPBAny 从protobuf Any 消息获取服务名
func TryServerNameFromPBAny(msg *any.Any) string {

	key, err := pb.AnyMessageNameTrimed(msg)
	if err != nil {
		panic(err)
	}

	key = key[strings.IndexByte(key, '_')+1:]
	return strings.ToLower(key)
}
