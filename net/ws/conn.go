package ws

import (
	"strconv"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/gorilla/websocket"
	"github.com/hiank/think/net/pb"
	"google.golang.org/protobuf/proto"
)

//conn 将websocket conn 的读写，转换为pool.ConnHandler 的读写
type conn struct {
	*websocket.Conn
	uid uint64
}

//GetKey 获取关键字
func (c *conn) Key() string {
	return strconv.FormatUint(c.uid, 10)
}

//Recv 读消息，实现frame.Conn
func (c *conn) Recv() (msg *pb.Message, err error) {

	_, buf, err := c.ReadMessage()
	if err == nil { //NOTE: 从websocket 读取消息
		a := new(any.Any)
		if err = proto.Unmarshal(buf, a); err == nil {
			msg = &pb.Message{SenderUid: c.uid, Value: a}
		}
	}
	return
}

//Send Writer
func (c *conn) Send(msg *pb.Message) (err error) {

	buf, err := proto.Marshal(msg.GetValue())
	if err == nil {
		err = c.WriteMessage(websocket.BinaryMessage, buf)
	}
	return
}
