package ws

import (
	"github.com/gorilla/websocket"
	"github.com/hiank/think/auth"
	"github.com/hiank/think/net"
	"github.com/hiank/think/run"
	"k8s.io/klog/v2"
)

const (
	ErrNotBinaryMessage = run.Err("ws: message recved should be BinaryMessage")
)

type conn struct {
	tk auth.Token
	wc *websocket.Conn
}

func (c *conn) Token() auth.Token {
	return c.tk
}

func (c *conn) Send(m *net.Message) error {
	return c.wc.WriteMessage(websocket.BinaryMessage, m.Bytes())
}

func (c *conn) Recv() (out *net.Message, err error) {
	mt, buf, err := c.wc.ReadMessage()
	if err == nil {
		switch mt {
		case websocket.BinaryMessage:
			defer func() {
				if r := recover(); r != nil {
					err = r.(error)
				}
			}()
			out = net.NewMessage(net.WithMessageBytes(buf))
		default:
			err = ErrNotBinaryMessage
			klog.Warning("ws: unsupport message type:", mt)
		}
	}
	return
}

func (c *conn) Close() error {
	return c.wc.Close()
}
