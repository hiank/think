package rpc

import (
	"context"

	"github.com/hiank/think/net"
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

func (c *conn) Send(b *net.Doc) (err error) {
	if err = c.ctx.Err(); err == nil {
		err = c.s.Send(b.Any())
	}
	return
}

func (c *conn) Recv() (out *net.Doc, err error) {
	if err = c.ctx.Err(); err == nil {
		var amsg *anypb.Any
		if amsg, err = c.s.Recv(); err == nil {
			out, err = net.MakeDoc(amsg)
		}
	}
	return
}

func (c *conn) Close() error {
	c.cancel()
	return nil
}
