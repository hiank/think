package ws

import (
	"github.com/gorilla/websocket"
	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/run"
)

const ErrUnsupportMessageType = run.Err("ws: unsupport message type recved")

type conn struct {
	// uid uint64
	wc *websocket.Conn
}

func (l *conn) Send(m pb.M) error {
	return l.wc.WriteMessage(websocket.BinaryMessage, m.Bytes())
}

func (l *conn) Recv() (out pb.M, err error) {
	t, bs, err := l.wc.ReadMessage()
	if err == nil {
		if t == websocket.BinaryMessage {
			out, err = pb.MakeM(bs)
		} else {
			err = ErrUnsupportMessageType
		}
	}
	return
}

func (c *conn) Close() error {
	return c.wc.Close()
}
