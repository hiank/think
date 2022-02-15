package ws

import (
	"github.com/gorilla/websocket"
	"github.com/hiank/think/net"
	"github.com/hiank/think/run"
)

const ErrUnsupportMessageType = run.Err("ws: unsupport message type recved")

type conn struct {
	// uid uint64
	wc *websocket.Conn
}

func (l *conn) Send(d *net.Doc) error {
	return l.wc.WriteMessage(websocket.BinaryMessage, d.Bytes())
}

func (l *conn) Recv() (out *net.Doc, err error) {
	t, bs, err := l.wc.ReadMessage()
	if err == nil {
		if t == websocket.BinaryMessage {
			out, err = net.MakeDoc(bs)
		} else {
			err = ErrUnsupportMessageType
		}
	}
	return
}

func (c *conn) Close() error {
	return c.wc.Close()
}
