package mq

import (
	"github.com/nats-io/nats.go"
)

//Client nats client connection
type Client struct {
	*nats.Conn
}

//TryNewClient 尝试创建一个Client，如果连接失败的话，会创建失败
func TryNewClient(url string) (c *Client, err error) {

	var nc *nats.Conn
	if nc, err = nats.Connect(url); err == nil {
		c = &Client{
			Conn: nc,
		}
	}
	return
}

// func (c *Client)
