package k8s

import (
	"strings"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/hiank/think/core"
	"github.com/hiank/think/core/pb"
)

//TryServerNameFromPBAny 从protobuf Any 消息获取服务名
func TryServerNameFromPBAny(msg *any.Any) string {

	key, err := pb.AnyMessageNameTrimed(msg)
	core.Panic(err)
	key = key[strings.IndexByte(key, '_')+1:]
	key = strings.ToLower(key)
	return key
}
