//Conn 在Listen 之前也可以接受Send 调用，但会返回一个error
//未Listen 之前 发送的消息将存于缓存 hub 中

package pool

import (
	"errors"

	"github.com/golang/glog"
	"github.com/hiank/think/pb"
	tk "github.com/hiank/think/token"
)

//Conn 用于管理连接
type Conn struct {
	*tk.Token             //NOTE：用于维护生命周期
	rw        IO          //NOTE: 读写Message
	hub       *MessageHub //NOTE: 待发送 MessageHub
}

//newConn 构建新的conn
func newConn(tok *tk.Token, rw IO) *Conn {

	c := &Conn{
		Token: tok,
		rw:    rw,
	}
	c.hub = NewMessageHub(tok, c)
	c.hub.DoActive() //NOTE: 不需要加锁等待，所有待处理的数据可以立即执行
	return c
}

//Listen 开启监听，每个conn 只有第一次调用，才生效
//一切正常的话，会阻塞在读消息通道中
func (c *Conn) Listen(readHandler MessageHandler) (err error) {

	var msg *pb.Message
L:
	for {
		select {
		case <-c.Done():
			err = errors.New("Conn's token Done")
			break L
		default:
			if msg, err = c.rw.Recv(); err != nil {
				c.Cancel() //NOTE: 此处调用，用于移除绑定的token
				break L
			}
			c.ResetTimer() //NOTE: 收到消息成功时，重置超时定时器
			if err := readHandler.Handle(NewMessage(msg, c.Derive())); err != nil {
				glog.Warning("Conn tokened "+c.ToString(), err)
			}
		}
	}
	return
}

//Handle 发送消息
func (c *Conn) Handle(msg *Message) error {

	select {
	case <-msg.Done():
		return errors.New("message's context was done") //NOTE: 要发送的消息绑定的context 关闭了
	case <-c.Done():
		return errors.New("Conn's context was done")
	default:
	}
	if err := c.rw.Send(msg.Message); err != nil { //NOTE: 发送失败，连接出了问题，退出[此处可能需要优化]
		c.Cancel()
		return err
	}
	c.ResetTimer() //NOTE: 发送成功的话，重置超时定时器
	return nil
}

//Send 发送消息，同步
func (c *Conn) Send(msg *Message) error {

	errChan := make(chan error)
	c.hub.PushWithBack(msg, errChan)
	return <-errChan
}

//IO 收发接口
type IO interface {
	Recv() (*pb.Message, error) //NOTE: 接收Message
	Send(*pb.Message) error     //NOTE: 发送Message
}

type connBuilder func() *Conn
