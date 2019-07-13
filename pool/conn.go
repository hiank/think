package pool

import (
	"context"
	"container/list"
	"github.com/hiank/think/pb"
)

//Conn pool中维护的Conn
type Conn struct {

	ConnHandler
	*list.Element
	Timer
	// context.Context

	Cancel context.CancelFunc
}

//NewDefaultConn 创建一个新的默认Conn
func NewDefaultConn(h ConnHandler) *Conn {

	c := NewConn(h)
	c.Timer = &DefaultTimer{}
	c.SetInterval(600)
	return c
}

//NewConn 新建一个Conn
func NewConn(h ConnHandler) *Conn {

	conn := &Conn {

		ConnHandler 	: h,
	}
	return conn
}

// //LoopRead 循环读取数据，送入传入的chan中
// func (conn *Conn) LoopRead(r chan *pb.Message) {

// 	L: for {

// 		msg, err := conn.ReadMessage()
// 		switch err {
// 		case nil: r <- msg
// 		default:
// 			close(r)
// 			break L
// 		}
// 	}
// }


//ConnHandler 数据读写接口
type ConnHandler interface {

	Identifier

	ReadMessage() (*pb.Message, error)		//NOTE: 读取Message
	WriteMessage(*pb.Message) error 		//NOTE: 写入Message
	// Close()									//NOTE: 关闭
	HandleContext(context.Context)			//NOTE: 处理Conn的Context
}

//IgnoreHandleContext 忽略处理Context
type IgnoreHandleContext int
//HandleContext 实现HandleContext 方法
func (ihc IgnoreHandleContext) HandleContext(context.Context) {}