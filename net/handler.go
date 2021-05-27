package net

import (
	"fmt"

	"github.com/hiank/think/net/pb"
	"google.golang.org/protobuf/proto"
)

//HandleFunc 函数形式的pool.Handler
//net package 中，所有用于传递的消息需是*pb.Message，此方法包含类型转换
type HandlerFunc func(*pb.Message) error

//Handle pool.Handler必要方法
func (hf HandlerFunc) Handle(m proto.Message) error {
	if pbMsg, ok := m.(*pb.Message); ok {
		return hf(pbMsg)
	}
	return fmt.Errorf("m %v not pb.Meesage", m)
}
