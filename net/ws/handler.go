package ws

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gorilla/websocket"
	"github.com/hiank/think/pb"
)

//Handler 将websocket conn 的读写，转换为pool.ConnHandler 的读写
type Handler struct {
	*websocket.Conn
	tokenStr string //NOTE: string token
}

//Recv 读消息，实现frame.Conn
func (c *Handler) Recv() (msg *pb.Message, err error) {

	_, buf, err := c.ReadMessage() //NOTE: 从websocket 读取消息
	if err != nil {
		return
	}

	glog.Infoln("ws conn read message :", buf)
	var a *any.Any
	if a, err = pb.AnyDecode(buf); err == nil {
		glog.Infoln("ws conn any decode :", a)
		msg = &pb.Message{Token: c.tokenStr, Data: a}
	}
	return
}

//Send Writer
func (c *Handler) Send(msg *pb.Message) (err error) {

	var buf []byte
	if buf, err = pb.AnyEncode(msg.GetData()); err != nil {
		glog.Warning(err)
		return
	}
	err = c.WriteMessage(websocket.BinaryMessage, buf)
	switch err {
	case nil:
	default: //NOTE:	处理错误
		glog.Warning(err)
	}
	return
}
