package adapter

import (
	"io"

	"github.com/hiank/think/net"
)

// type ChanAccepter struct {
// 	pp chan net.IdentityConn
// }

// func NewChanAccepter() *ChanAccepter {
// 	return &ChanAccepter{
// 		pp: make(chan net.IdentityConn),
// 	}
// }

// func (ca *ChanAccepter) Chan() chan<- net.IdentityConn {
// 	return ca.pp
// }

// func (ca *ChanAccepter) Accept() (ic net.IdentityConn, err error) {
// 	ic, ok := <-ca.pp
// 	if !ok {
// 		err = io.EOF
// 	}
// 	return
// }

type ChanAccepter chan net.IdentityConn

func (ca ChanAccepter) Accept() (ic net.IdentityConn, err error) {
	ic, ok := <-ca
	if !ok {
		err = io.EOF
	}
	return
}
