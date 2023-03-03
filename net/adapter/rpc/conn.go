package rpc

import (
	"io"
	"time"

	"github.com/hiank/think/auth"
	"github.com/hiank/think/net"
	"github.com/hiank/think/run"
	"google.golang.org/protobuf/types/known/anypb"
)

type sendReciver interface {
	Send(*anypb.Any) error
	Recv() (*anypb.Any, error)
}

type conn struct {
	tk auth.Token
	sr sendReciver
	io.Closer
}

func (c *conn) Token() auth.Token {
	return c.tk
}

func (c *conn) Send(m *net.Message) (err error) {
	if err = run.FrontErr(m.Token().Err, c.tk.Err); err == nil {
		err = c.sr.Send(m.Any())
	}
	return
}

func (c *conn) Recv() (out *net.Message, err error) {
	if err = c.tk.Err(); err == nil {
		var amsg *anypb.Any
		if amsg, err = c.sr.Recv(); err == nil {
			out = net.NewMessage(net.WithMessageValue(amsg), net.WithMessageToken(c.tk.Fork(auth.WithTokenTimeout(time.Second*5))))
		}
	}
	return
}

// type restClient struct {
// 	pipe.RestClient
// 	io.Closer
// }
