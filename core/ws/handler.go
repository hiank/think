package ws

import (
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gorilla/websocket"
	"github.com/hiank/think/core"
	"github.com/hiank/think/core/pb"
	"google.golang.org/protobuf/proto"
)

//Handler 将websocket conn 的读写，转换为pool.ConnHandler 的读写
type Handler struct {
	*websocket.Conn
	tokenStr string //NOTE: string token
}

//Recv 读消息，实现frame.Conn
func (c *Handler) Recv() (msg *pb.Message, err error) {

	defer core.Recover(core.Warning)

	_, buf, err := c.ReadMessage() //NOTE: 从websocket 读取消息
	core.Panic(err)

	a := new(any.Any)
	err = proto.Unmarshal(buf, a)
	core.Panic(err)

	msg = &pb.Message{Key: c.tokenStr, Value: a}
	return
}

//Send Writer
func (c *Handler) Send(msg *pb.Message) (err error) {

	defer core.Recover(core.Warning)

	buf, err := proto.Marshal(msg.GetValue())
	core.Panic(err)

	err = c.WriteMessage(websocket.BinaryMessage, buf)
	core.Panic(err)
	return
}
