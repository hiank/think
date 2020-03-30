package pool

import (
	"context"
	"errors"
	"sync"

	"github.com/golang/glog"
	"github.com/hiank/think/pb"
)

//ConnHub 用于存储管理Conn
type ConnHub struct {

	ctx 	context.Context				//NOTE: 
	mtx 	sync.RWMutex				//NOTE: 读写锁，ConnHub 会在不同goroutine中添删conn
	// queue 	*list.List					//NOTE: 用于保存连接
	hub 	map[string]*Conn			//NOTE: map[tokenString]*Conn
}


//NewConnHub 构建ConnHub
func NewConnHub(ctx context.Context) *ConnHub {

	ch := &ConnHub{
		ctx 	: ctx,
		hub 	: make(map[string]*Conn),
	}
	return ch
}


//AddConn 添加新的Conn，启用相关处理goroutine
func (ch *ConnHub) AddConn(conn *Conn) (err error) {

	ch.push(conn)
	L: for err == nil {
		select {
		case <-conn.Done():
			err = errors.New("conn tokened : " + conn.ToString() + " Done")
			break L
		case <-ch.ctx.Done():				//NOTE: Context 被关闭了，执行清理
			err = errors.New("conn tokened : " + conn.ToString() + " Done")
			break L
		default:
			err = ch.read(conn)
		}
	}
	ch.Remove(conn)
	return
}

func (ch *ConnHub) read(conn *Conn) error {

	msg, err := conn.Recv()
	select {
	case <-ch.ctx.Done():
		return errors.New("conn tokened : " + conn.ToString() + " Done")
	case <-conn.Done():
		return errors.New("conn tokened : " + conn.ToString() + " Done")
	default:
	}		//NOTE: 此处确保连接关闭后，不再处理后续

	switch err {
	case nil:
		ch.update(conn)
		go ch.ctx.Value(CtxKeyRecvHandler).(MessageHandler).Handle(msg)		//NOTE: 处理收到的消息
	default:
		glog.Warningln("conn read error : ", err, "...tokened : ", conn.ToString())
		conn.Cancel()		//NOTE: 连接出错，释放token，所有相关的资源将被释放
	}
	return err
}


func (ch *ConnHub) push(conn *Conn) (err error) {

	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	ch.hub[conn.ToString()] = conn
	// conn.Element = ch.queue.PushBack(conn)
	// ch.queue.PushBack(conn)
	return
}

func (ch *ConnHub) update(conn *Conn) {

	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	select {
	case <-conn.Done():		//NOTE: 极端情况下，upgrade 将此conn Cancel了，当锁解开后，此处conn 已经失效，所以加个判断
	default:
		conn.ResetTimer()	//NOTE: 重置相关token 的定时器
		// ch.queue.MoveToBack(conn.Element)
		// conn.Reset()
	}
}


//Handle 处理数据发送，总感觉这边会有性能问题，如果有超级多的玩家同时在线，比如1000万，每次要发送一个消息都要遍历查找一遍，可能会卡死
func (ch *ConnHub) Handle(msg *pb.Message) error {

	ch.mtx.RLock()
	defer ch.mtx.RUnlock()

	if conn, ok := ch.hub[msg.GetToken()]; ok {
		return conn.Send(msg)
	}
	if builder := ch.ctx.Value(CtxKeyConnBuilder).(ConnBuilder); builder != nil {
		go builder.BuildAndSend(msg)
	}
	return errors.New("connhub has no conn tokened " + msg.GetToken())
}

//Remove 删除Conn
func (ch *ConnHub) Remove(conn *Conn) (err error) {

	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	delete(ch.hub, conn.ToString())
	return
}


//ConnBuilder interface for build Conn
type ConnBuilder interface {

	BuildAndSend(msg *pb.Message)
}