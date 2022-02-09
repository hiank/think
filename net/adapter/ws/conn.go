package ws

import (
	"errors"

	"github.com/gorilla/websocket"
	"github.com/hiank/think/net"
)

type conn struct {
	// uid uint64
	wc *websocket.Conn
}

// func (c *conn) GetIdentity() uint64 {
// 	return c.uid
// }

func (l *conn) Send(d *net.Doc) error {
	return l.wc.WriteMessage(websocket.BinaryMessage, d.Bytes())
}

func (l *conn) Recv() (out *net.Doc, err error) {
	t, bs, err := l.wc.ReadMessage()
	if t != websocket.BinaryMessage {
		err = errors.New("only support 'BinaryMessage' type (protobuf)")
	} else if err == nil {
		out, err = net.MakeDoc(bs)
	}
	return
}

func (c *conn) Close() error {
	return c.wc.Close()
}
