package net

import (
	"context"
	"io"

	"github.com/hiank/think/net/box"
	"google.golang.org/protobuf/proto"
)

type Sender interface {
	Send(*box.Message) error
}

type Receiver interface {
	Recv() (*box.Message, error)
}

type Conn interface {
	Sender
	Receiver
	io.Closer
}

// type Rest interface {
// 	Get(context.Context, *anypb.Any) (*anypb.Any, error)
// 	Post(context.Context, *anypb.Any) (*emptypb.Empty, error)
// 	io.Closer
// }

//IdentityConn conn with identity
type IdentityConn struct {
	ID string
	Conn
}

type IdentityMessage struct {
	ID string
	M  *box.Message
}

//Knower unpack box.Message to want data
type Knower interface {
	//Identiy uid to identity
	Identity(uid string) (string, error)

	//ServeAddr get server addr from box.message
	ServeAddr(*box.Message) (string, error)
}

//Dialer
type Dialer interface {
	Dial(ctx context.Context, addr string) (Conn, error)
}

type Listener interface {
	Accept() (IdentityConn, error)
	Close() error
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
	Route(string, *box.Message)
}

//HandlerFunc easy convert func to Handler
type HandlerFunc func(id string, m *box.Message)

func (hf HandlerFunc) Route(id string, m *box.Message) {
	hf(id, m)
}
