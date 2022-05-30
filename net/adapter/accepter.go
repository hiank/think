package adapter

import (
	"io"

	"github.com/hiank/think/net"
)

type ChanAccepter chan net.TokenConn

func (ca ChanAccepter) Accept() (ic net.TokenConn, err error) {
	ic, ok := <-ca
	if !ok {
		err = io.EOF
	}
	return
}
