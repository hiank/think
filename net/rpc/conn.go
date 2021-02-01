package rpc

import (
	"io"

	"github.com/hiank/think/net"
)

//Conn grpc conn
type Conn struct {
	net.Sender
	net.Reciver
	io.Closer
	key string
}

//Key 关键字
func (c *Conn) Key() string {
	return c.key
}

//Close 关闭连接
func (c *Conn) Close() (err error) {

	if c.Closer != nil {
		err = c.Closer.Close()
	}
	return
}
