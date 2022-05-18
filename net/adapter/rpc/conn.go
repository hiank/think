package rpc

import (
	"context"
	"io"

	"github.com/hiank/think/net/adapter/rpc/pipe"
	"github.com/hiank/think/net/box"
	"google.golang.org/protobuf/types/known/anypb"
)

type conn struct {
	ctx    context.Context
	cancel context.CancelFunc
	s      SendReciver
}

func (c *conn) Send(m *box.Message) (err error) {
	if err = c.ctx.Err(); err == nil {
		err = c.s.Send(m.GetAny())
	}
	return
}

func (c *conn) Recv() (out *box.Message, err error) {
	if err = c.ctx.Err(); err == nil {
		var amsg *anypb.Any
		if amsg, err = c.s.Recv(); err == nil {
			out, err = box.New(amsg)
		}
	}
	return
}

func (c *conn) Close() error {
	c.cancel()
	return nil
}

type restClient struct {
	pipe.RestClient
	io.Closer
}
