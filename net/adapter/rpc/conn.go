package rpc

import (
	"context"

	"github.com/hiank/think/net/pb"
	"google.golang.org/protobuf/types/known/anypb"
)

// type iSendRecver interface {
// 	Send(*pb.Carrier) error
// 	Recv() (*pb.Carrier, error)
// }

type conn struct {
	ctx    context.Context
	cancel context.CancelFunc
	s      Stream
	// identity uint64
}

// func (c *conn) GetIdentity() uint64 {
// 	return c.identity
// }

func (c *conn) Send(m pb.M) (err error) {
	if err = c.ctx.Err(); err == nil {
		err = c.s.Send(m.Any())
	}
	return
}

func (c *conn) Recv() (out pb.M, err error) {
	if err = c.ctx.Err(); err == nil {
		var amsg *anypb.Any
		if amsg, err = c.s.Recv(); err == nil {
			out, err = pb.MakeM(amsg)
		}
	}
	return
}

func (c *conn) Close() error {
	c.cancel()
	return nil
}
