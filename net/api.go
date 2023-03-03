package net

import (
	"context"
	"io"

	"github.com/hiank/think/auth"
)

type Sender interface {
	Send(*Message) error
}

type Receiver interface {
	Recv() (*Message, error)
}

type Conn interface {
	Token() auth.Token
	Sender
	Receiver
	io.Closer
}

// Knower unpack Message to want data
type Knower interface {
	//ServeAddr get server addr from message
	ServeAddr(*Message) (string, error)
}

// Dialer
type Dialer interface {
	Dial(ctx context.Context, addr string) (Conn, error)
}

// Listner for server 
type Listener interface {
	Accept() (Conn, error)
	Close() error

}
// Handler handle message
type Handler interface {
	Route(*Message)
}

// HandlerFunc easy convert func to Handler
type FuncHandler func(msg *Message)

func (fh FuncHandler) Route(msg *Message) {
	fh(msg)
}
