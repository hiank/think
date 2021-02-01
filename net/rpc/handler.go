package rpc

import (
	"github.com/hiank/think/net/pb"
)

//*****************************handle message recv from grpc*******************************//

// //ReadHandler 服务器端消息处理
// //处理服务端收到的消息
// type ReadHandler interface {
// 	net.MessageHandler //NOTE: 处理stream消息
// }

//DonceHandler 'Get' 'Post'类型消息的处理接口
type DonceHandler interface {
	HandleGet(*pb.Message) (*pb.Message, error) //NOTE: 处理Get消息
	HandlePost(*pb.Message) error               //NOTE: 处理Post消息
}

//IgnoreGet 忽略Get 实现
type IgnoreGet int

//HandleGet 用于忽略HandleGet
func (i IgnoreGet) HandleGet(*pb.Message) (msg *pb.Message, err error) { return }

//IgnorePost 忽略Post 实现
type IgnorePost int

//HandlePost 用于忽略HandlePost 方法
func (i IgnorePost) HandlePost(*pb.Message) (err error) { return }

// //IgnoreStream 忽略Stream 实现
// type IgnoreStream int

// //Handle 用于忽略Handle
// func (i IgnoreStream) Handle(core.Message) (err error) { return }
