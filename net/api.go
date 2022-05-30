package net

import (
	"context"
	"io"

	"github.com/hiank/think/net/box"
	"google.golang.org/protobuf/proto"
)

type Sender interface {
	Send(box.Message) error
}

type Receiver interface {
	Recv() (box.Message, error)
}

type Conn interface {
	Sender
	Receiver
	io.Closer
}

//TokenConn struct contianed Token and Conn
type TokenConn box.TT[Conn]

//TokenMessage struct contianed Token and *Message
type TokenMessage box.TT[box.Message]

//Knower unpack box.Message to want data
type Knower interface {
	//Identiy uid to identity
	// Identity(uid string) (string, error)

	//ServeAddr get server addr from box.message
	ServeAddr(box.Message) (string, error)
}

//Dialer
type Dialer interface {
	Dial(ctx context.Context, addr string) (Conn, error)
}

type Listener interface {
	Accept() (TokenConn, error)
	Close() error
}

type Clientset interface {
	AutoSend(tm TokenMessage) error
	RouteMux() *RouteMux
	io.Closer
}

type Server interface {
	//start work
	ListenAndServe() error
	//Send message to client (by conn)
	//ti is target identity. when len(ti) == 0
	//means send for all conn
	Send(v proto.Message, tis ...string) error
	io.Closer
}

//Handler handle message
type Handler interface {
	// Route(string, *box.Message)
	Route(TokenMessage)
}

//HandlerFunc easy convert func to Handler
type FuncHandler func(tt TokenMessage)

func (fh FuncHandler) Route(tt TokenMessage) {
	fh(tt)
}
