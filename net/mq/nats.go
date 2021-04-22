package mq

import (
	"github.com/nats-io/nats.go"
)

//NewNatsConn 创建一个nats连接
func NewNatsConn(url string) (*nats.Conn, error) {
	return nats.Connect(url)
}
