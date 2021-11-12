package rpc

import (
	"github.com/hiank/think/net"
)

type dialer struct {
}

func NewDialer() net.IDialer {
	d := &dialer{}
	return d
}

func (d *dialer) Dial(addr string) (net.IConn, error) {
	return nil, nil
}
