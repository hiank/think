package ws

import (
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gorilla/websocket"
	"github.com/hiank/think/pb"
	"github.com/hiank/think/utils/robust"
	"google.golang.org/protobuf/proto"
)

//Handler 将websocket conn 的读写，转换为pool.ConnHandler 的读写
type Handler struct {
	*websocket.Conn
	tokenStr string //NOTE: string token
}

//Recv 读消息，实现frame.Conn
func (c *Handler) Recv() (msg *pb.Message, err error) {

	defer robust.Recover(robust.Warning)

	_, buf, err := c.ReadMessage() //NOTE: 从websocket 读取消息
	robust.Panic(err)

	a := new(any.Any)
	err = proto.Unmarshal(buf, a)
	robust.Panic(err)

	msg = &pb.Message{Token: c.tokenStr, Data: a}
	return
}

//Send Writer
func (c *Handler) Send(msg *pb.Message) (err error) {

	defer robust.Recover(robust.Warning)

	buf, err := proto.Marshal(msg.GetData())
	robust.Panic(err)

	err = c.WriteMessage(websocket.BinaryMessage, buf)
	robust.Panic(err)
	return
}
