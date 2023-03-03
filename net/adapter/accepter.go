package adapter

import (
	"io"

	"github.com/hiank/think/net"
)

type ChanAccepter chan net.Conn

func (ca ChanAccepter) Accept() (ic net.Conn, err error) {
	ic, ok := <-ca
	if !ok {
		err = io.EOF
	}
	return
}
