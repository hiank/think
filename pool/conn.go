//Conn 在Listen 之前也可以接受Send 调用，但会返回一个error
//未Listen 之前 发送的消息将存于缓存 hub 中

package pool

import (
	"errors"

	"github.com/golang/glog"
	"github.com/hiank/think/pb"
	tk "github.com/hiank/think/token"
)

//conn 用于管理连接
type conn struct {
	*tk.Token   		//NOTE：用于维护生命周期
	rw 		IO 				//NOTE: 读写Message
	hub 	*MessageHub		//NOTE: 待发送 MessageHub

	exit 	chan error		//NOTE: 退出指令，当读或写消息出错后，通知此chan，用于退出Listen，并结束Conn
}

//newConn 构建新的conn
func newConn(tok *tk.Token, rw IO) *conn {

	c := &conn{
		Token: 	tok,
		rw: 	rw,
		exit: 	make(chan error),
	}
	c.hub = NewMessageHub(tok, c)
	return c
}


//Listen 开启监听，每个conn 只有第一次调用，才生效
//一切正常的话，会阻塞在读消息通道中
func (c *conn) Listen(readHandler MessageHandler) error {

	select {
	case <-c.Done(): return errors.New("conn tokend " + c.ToString() + " Done")
	default:
	}

	go c.loopRead(readHandler)			//NOTE: 起一个读协程
	err := <- c.exit
	c.Cancel()
	return err
}


//Handle 发送消息
func (c *conn) Handle(msg *Message) error {

	select {
	case <-msg.Done(): return errors.New("message's context was done")		//NOTE: 要发送的消息绑定的context 关闭了
	case <-c.Done(): return errors.New("conn's context was done")
	default:
	}
	if err := c.rw.Send(msg.Message); err != nil {		//NOTE: 发送失败，人物连接出了问题，退出[此处可能需要优化]
		c.exit <- err
		return err
	}
	c.ResetTimer()								//NOTE: 发送成功的话，重置超时定时器
	return nil
}


//Send 发送消息，同步
func (c *conn) Send(msg *Message) error {

	errChan := make(chan error)
	c.hub.PushWithBack(msg, errChan)
	return <-errChan
}


//loopRead 循环读消息
func (c *conn) loopRead(handler MessageHandler) {

	L: for {

		select {
		case <-c.Done(): 
			c.exit <- errors.New("conn's token Done")
			break L
		default:
		}
		msg, err := c.rw.Recv()
		if err != nil {
			c.exit <- err
			break L
		}
		c.ResetTimer()				//NOTE: 收到消息成功时，重置超时定时器
		if err = handler.Handle(NewMessage(msg, c.Derive())); err != nil {
			glog.Warning("conn tokend " + c.ToString(), err)
		}
	}
}
 

//IO 收发接口
type IO interface {
	Recv() (*pb.Message, error) 	//NOTE: 接收Message
	Send(*pb.Message) error     	//NOTE: 发送Message
}

