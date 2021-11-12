package ws

import (
	"errors"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type conn struct {
	uid uint64
	wc  *websocket.Conn
}

func (c *conn) GetIdentity() uint64 {
	return c.uid
}

func (l *conn) Send(any *anypb.Any) error {
	b, err := proto.Marshal(any)
	if err == nil {
		err = l.wc.WriteMessage(websocket.BinaryMessage, b)
	}
	return err
}

func (l *conn) Recv() (any *anypb.Any, err error) {
	t, b, err := l.wc.ReadMessage()
	if t != websocket.BinaryMessage {
		err = errors.New("only support 'BinaryMessage' type (protobuf)")
	} else if err == nil {
		any = &anypb.Any{}
		err = proto.Unmarshal(b, any)
	}
	return
}

func (c *conn) Close() error {
	return c.wc.Close()
}
