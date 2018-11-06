package protobuf

import (
	"strings"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"

	tp "github.com/hiank/think/net/protobuf/think"
)



// PBDecode used to unmarshal net data to expect data
func PBDecode (msg []byte) (*any.Any, error) {

	bus := &tp.BUS{}
	err := proto.Unmarshal(msg, bus)
	if err != nil {
		return nil, err
	}
	return bus.GetData(), nil
}

// PBEncode used to marshal data to net data
func PBEncode (anyMsg *any.Any) ([]byte, error) {

	bus := &tp.BUS{Data:anyMsg,}
	buf, err := proto.Marshal(bus)
	if err != nil {
		return nil, err
	}
	return buf, nil
}


// GetServerName 通过message name 获得服务名
func GetServerName (messageName string) string {

	idx := strings.IndexByte(messageName, '_')
	serverName := messageName[:idx]
	return serverName
}

