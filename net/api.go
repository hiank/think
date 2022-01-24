package net

import (
	"github.com/hiank/think/net/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Conn interface {
	GetIdentity() uint64
	Recv() (*anypb.Any, error)
	Send(*anypb.Any) error
	Close() error
}

type Dialer interface {
	Dial(addr string) (Conn, error)
}

type Client interface {
	Send(*pb.Carrier) error
	Close() error
}

type Listener interface {
	Accept() (Conn, error)
	Close() error
}

type Server interface {
	ListenAndServe() error
	//Send message to client (by conn)
	//NOTE: the method will wait until completed
	Send(*pb.Carrier) error
	Close() error
}

type CarrierHandler interface {
	Handle(*pb.Carrier)
}

type MessageHandler interface {
	Handle(id uint64, msg proto.Message)
}

//CarrierConverter convert Carrier to string key (use in HandleMux)
//HandleMux use the converter to known which CarrierHandler registered use to Handle carrier message
type CarrierConverter interface {
	GetKey(*pb.Carrier) (key string, ok bool)
}

var (
	//NewServer new a Server.
	//CarrierHandler use to handle message received from client (by Conn)
	NewServer func(Listener, CarrierHandler) Server = newServer

	NewClient func(Dialer) Client = newClient
)
