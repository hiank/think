package core

import (
	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/ptypes/any"
)

//RedisHolder get redis
type RedisHolder interface {
	TryMaster() *redis.Client
	TrySlave() *redis.Client
}

//Message 消息
type Message interface {
	GetKey() string //NOTE: 消息关键字，与Conn及Service 的key对应
	GetValue() *any.Any
}

//MessageHandler 消息处理者
type MessageHandler interface {
	Handle(Message) error
}

//Conn 连接句柄
type Conn interface {
	GetKey() string         //NOTE: 消息关键字，与Conn及Service 的key对应
	Recv() (Message, error) //NOTE: 收消息
	Send(Message) error     //NOTE: 发消息
	Close() error           //NOTE: 关闭
}
