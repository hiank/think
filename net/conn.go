package net

import (
	"io"

	"github.com/hiank/think/net/pb"
)

//Conn 消息读写接口
type Conn interface {
	Key() string //NOTE: 连接关键字
	Sender
	Reciver
	io.Closer
}

//Sender 发送接口
type Sender interface {
	Send(*pb.Message) error
}

//Reciver 接收接口
type Reciver interface {
	Recv() (*pb.Message, error)
}

//SenderFunc 函数形式Sender
type SenderFunc func(*pb.Message) error

//Send sender实现
func (sf SenderFunc) Send(msg *pb.Message) error {
	return sf(msg)
}

//ReciverFunc 函数形式Reciver
type ReciverFunc func() (*pb.Message, error)

//Recv reciver实现
func (rf ReciverFunc) Recv() (*pb.Message, error) {
	return rf()
}

//CloserFunc 函数形式io.Closer
type CloserFunc func() error

//Close io.Close实现
func (cf CloserFunc) Close() error {
	return cf()
}
