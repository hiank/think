package net

import (
	"github.com/hiank/think/net/pb"
	"google.golang.org/protobuf/types/known/anypb"
)

type IConn interface {
	GetIdentity() uint64
	Recv() (*anypb.Any, error)
	Send(*anypb.Any) error
	Close() error
}

type IDialer interface {
	Dial(addr string) (IConn, error)
}

type IClient interface {
	Send(*pb.Carrier) error
	Close() error
}

type IListener interface {
	Accept() (IConn, error)
	Close() error
}

type IServer interface {
	ListenAndServe() error
	//Send message to client (by conn)
	//NOTE: the method will wait until completed
	Send(*pb.Carrier) error
	Close() error
}

type ICarrierHandler interface {
	Handle(*pb.Carrier)
}

//ICarrierConverter convert Carrier to string key (use in HandleMux)
//HandleMux use the converter to known which ICarrierHandler registered use to Handle carrier message
type ICarrierConverter interface {
	GetKey(*pb.Carrier) (key string, ok bool)
}

var (
	//NewServer new a IServer.
	//ICarrierHandler use to handle message received from client (by IConn)
	NewServer func(IListener, ICarrierHandler) IServer = newServer

	NewClient func(IDialer) IClient = newClient
)
