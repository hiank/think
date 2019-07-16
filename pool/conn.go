package pool

import (
	"context"
	"github.com/golang/glog"
	"github.com/hiank/think/token"
	"errors"
	"container/list"
	"github.com/hiank/think/pb"
)

//Conn pool中维护的Conn
type Conn struct {

	*list.Element
	Timer

	key 	string			//NOTE: 
	tk 		*token.Token
	ch 		ConnHandler
}

//newDefaultConn 创建一个新的默认Conn
func newDefaultConn(key string, t *token.Token, h ConnHandler) *Conn {

	conn := &Conn {

		tk 		: t,
		ch 		: h,
		key 	: key,
	}
	conn.Timer = &DefaultTimer{}
	conn.SetInterval(600)
	return conn
}

//NewConn 新建一个Conn，会用到tokenStr绑定的Token，这个Token 出现异常导致Done的情况下，所有tokenStr 绑定的Token 都会Done
func NewConn(key, tokenStr string, h ConnHandler) *Conn {

	t, err := token.Get(tokenStr)
	if err != nil {
		glog.Warningln("NewConn error : ", err)
		return nil
	}
	return newDefaultConn(key, t, h)
}

//NewConnWithDerivedToken 使用派生Token 生成的Conn，如果与grpc服务连接的Conn，生命周期独立，连接异常断开 不会影响到其它使用了这个tokenStr 的代码
func NewConnWithDerivedToken(key, tokenStr string, h ConnHandler) *Conn {

	t, err := token.Get(tokenStr)
	if err != nil {
		glog.Warningln("NewConn error : ", err)
		return nil
	}
	return newDefaultConn(key, t, h)
}

//GetKey 获得Conn关键字，用于分类
func (conn *Conn) GetKey() string {

	return conn.key
}

//GetToken 获得conn 的Token
func (conn *Conn) GetToken() *token.Token {

	return conn.tk
}

//Send 发送消息
func (conn *Conn) Send(msg *pb.Message) (err error) {

	select {
	case <-conn.tk.Done():		//NOTE: 这个Token 被关闭了[或者是整个应用被关闭了]
		err = errors.New("Token " + conn.tk.ToString() + " cancelled")
	default:
		err = conn.ch.WriteMessage(conn.GetToken().WithValue("key", conn.GetKey()), msg)
	}
	return
}

//Recv 接收消息
func (conn *Conn) Recv() (msg *pb.Message, err error) {

	select {
	case <-conn.tk.Done():		//NOTE: 这个Token 被关闭了[或者是整个应用被关闭了]
		err = errors.New("Token " + conn.tk.ToString() + " cancelled")
	default:
		msg, err = conn.ch.ReadMessage(conn.GetToken().WithValue("key", conn.GetKey()))
	}
	return
}


//ConnHandler 数据读写接口
type ConnHandler interface {

	ReadMessage(ctxWithKeyToken context.Context) (*pb.Message, error)			//NOTE: 读取Message，传入key token信息，用于构建或处理msg
	WriteMessage(ctxWithKeyToken context.Context, msg *pb.Message) error 		//NOTE: 写入Message，传入key token信息，用于构建货处理msg
}


// //IgnoreHandleContext 忽略处理Context
// type IgnoreHandleContext int
// //HandleContext 实现HandleContext 方法
// func (ihc IgnoreHandleContext) HandleContext(context.Context) {}