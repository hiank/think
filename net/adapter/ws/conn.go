package ws

import (
	"github.com/gorilla/websocket"
	"github.com/hiank/think/net/box"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/klog/v2"
)

type conn struct {
	wc *websocket.Conn
}

func (c *conn) Send(m box.Message) error {
	return c.wc.WriteMessage(websocket.BinaryMessage, m.GetBytes())
}

func (c *conn) Recv() (out box.Message, err error) {
	mt, buf, err := c.wc.ReadMessage()
	if err == nil {
		switch mt {
		case websocket.BinaryMessage:
			out = box.New() //new(box.Message)
			err = box.Unmarshal[*anypb.Any](buf, out)
		default:
			klog.Warning("ws: unsupport message type:", mt)
		}
	}
	return
}

func (c *conn) Close() error {
	return c.wc.Close()
}
