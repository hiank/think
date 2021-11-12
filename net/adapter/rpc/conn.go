package rpc

import (
	"context"

	"github.com/hiank/think/net/pb"
	"google.golang.org/protobuf/types/known/anypb"
)

type iSendRecver interface {
	Send(*pb.Carrier) error
	Recv() (*pb.Carrier, error)
}

type conn struct {
	ctx      context.Context
	cancel   context.CancelFunc
	sr       iSendRecver
	identity uint64
}

func (c *conn) GetIdentity() uint64 {
	return c.identity
}

func (c *conn) Send(any *anypb.Any) error {
	return c.sr.Send(&pb.Carrier{Identity: c.identity, Message: any})
	// return nil
}

func (c *conn) Recv() (any *anypb.Any, err error) {
	carrier, err := c.sr.Recv()
	if err == nil {
		any = carrier.GetMessage()
	}
	return
}

func (c *conn) Close() error {
	c.cancel()
	return nil
}
