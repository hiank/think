package net

import "github.com/hiank/think/net/pb"

type Conn interface {
	Send(pb.M) error
	Recv() (pb.M, error)
	Close() error
}

//IAC identity and conn
type IAC struct {
	ID string
	Conn
}

//Dialer dial to server
type Dialer interface {
	Dial(addr string) (IAC, error)
}

type Client interface {
	Send(d pb.M, ti string) error
	// Handle(k any, h Handler)
}

type Listener interface {
	Accept() (IAC, error)
	Close() error
}

type Server interface {
	//start work
	ListenAndServe() error
	//Send message to client (by conn)
	//ti is target identity. when len(ti) == 0
	//means send for all conn
	Send(v any, tis ...string) error
	// //AddHandler add handler for revc message
	// Handle(k any, h Handler)
	//
	Close() error
}

//Handler handle message
type Handler interface {
	Route(string, pb.M)
}

//HandlerFunc easy convert func to Handler
type HandlerFunc func(string, pb.M)

func (hf HandlerFunc) Route(id string, m pb.M) {
	hf(id, m)
}
