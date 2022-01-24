package rpc

import (
	"github.com/hiank/think/net"
)

type dialer struct {
}

func NewDialer() net.Dialer {
	d := &dialer{}
	return d
}

func (d *dialer) Dial(addr string) (net.Conn, error) {
	return nil, nil
}
