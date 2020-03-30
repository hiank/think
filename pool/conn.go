package pool

import (
	tk "github.com/hiank/think/token"
	"github.com/hiank/think/pb"
)

//Conn 用于管理连接
type Conn struct {
	*tk.Token   	//NOTE：用于维护生命周期
	IO 				//NOTE: 读写Message
}

//NewConn 构建新的Conn
func NewConn(tok *tk.Token, rw IO) *Conn {

	return &Conn{
		Token: 	tok,
		IO: 	rw,
	}
}

//IO 收发接口
type IO interface {
	Recv() (*pb.Message, error) //NOTE: 接收Message
	Send(*pb.Message) error     //NOTE: 发送Message
}
