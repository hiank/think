package pb

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/proto"

	"strings"
	"github.com/golang/protobuf/ptypes/any"
)


// // Message2Any protobuf message to Any
// func Message2Any(msg proto.Message) (*any.Any, error) {

// 	buf, err := proto.Marshal(msg)
// 	if err != nil {
// 		return nil, err
// 	}
	
// 	anyMsg, err := AnyDecode(buf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return anyMsg, nil
// }


// AnyDecode used to unmarshal net data to expect data
func AnyDecode(msg []byte) (*any.Any, error) {

	a := &any.Any{}
	err := proto.Unmarshal(msg, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// AnyEncode used to marshal data to net data
func AnyEncode(anyMsg *any.Any) ([]byte, error) {

	buf, err := proto.Marshal(anyMsg)
	if err != nil {
		return nil, err
	}
	return buf, nil
}


// GetServerKey 通过message name 获得服务名
func GetServerKey(anyMsg *any.Any) (n string, e error) {

	messageName, e := ptypes.AnyMessageName(anyMsg)
	if e != nil {

		glog.Warningf("get any message name error : %v\n", e)
		return
	}

	if dotIdx := strings.LastIndexByte(messageName, '.'); dotIdx != -1 {
		messageName = messageName[dotIdx+1:]		//NOTE: 协议可能包含包名，此处截掉包名
	}

	glog.Infof("messageName : %s\n", messageName)
	messageName = messageName[2:]					//NOTE: 前两位用于保存消息类型
	idx := strings.IndexByte(messageName, '_')
	n = strings.ToLower(messageName[:idx])
	return
}

//Message Type 用于甄别message 需要用那种方式调用
const (

	TypeUndefined 		= iota		//NOTE: 未定义的类型，表明出错了
	TypeGET 						//NOTE: get消息，需要一个返回
	TypePOST 						//NOTE: post消息
	TypeSTREAM 						//NOTE: 流消息
)

// GetServerType 获得服务类型
func GetServerType(anyMsg *any.Any) (t int, err error) {

	messageName, err := ptypes.AnyMessageName(anyMsg);
	if err != nil {
		return TypeUndefined, err
	}

	if dotIdx := strings.LastIndexByte(messageName, '.'); dotIdx != -1 {
		messageName = messageName[dotIdx+1:]		//NOTE: 协议可能包含包名，此处截掉包名
	}

	switch messageName[0] {
	case 'G': t = TypeGET
	case 'P': t = TypePOST
	case 'S': t = TypeSTREAM
	default: t = TypeUndefined
	}
	return
}
