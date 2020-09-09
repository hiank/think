package mq

import (
	"github.com/hiank/think/utils/robust"
	"github.com/nats-io/nats.go"
)

//Client nats client connection
type Client struct {
	*nats.Conn
}

//TryNewClient 尝试创建一个Client，如果连接失败的话，抛出异常，注意处理异常
func TryNewClient(url string) *Client {

	nc, err := nats.Connect(url)
	robust.Panic(err)
	return &Client{
		Conn: nc,
	}
}

// func (c *Client)
