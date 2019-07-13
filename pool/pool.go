package pool

import (
	"github.com/golang/glog"
	"errors"
	"github.com/hiank/think/pb"
	"github.com/hiank/think/conf"
	"time"
	"context"
	"sync"
)

//Pool 消息处理的核心
type Pool struct {

	sync.RWMutex
	*ConnHub

	readhub 	*MessageHub 			//NOTE: 读消息hub
	writehub 	*MessageHub 			//NOTE: 写消息hub

	ctx			context.Context	 		//NOTE: Context 用于维护生命周期
	Close 		context.CancelFunc 		//NOTE: 关闭方法
}


//NewPool 创建一个新的Pool
func NewPool(ctx context.Context, mh MessageHandler) *Pool {

	ctx, cancel := context.WithCancel(ctx)

	num := conf.GetSys().MaxMessageGo
	ch := NewConnHub()
	pool := &Pool {
		readhub 	: NewMessageHub(mh, num),
		writehub 	: NewMessageHub(ch, num),
		ConnHub 	: ch,
		ctx 		: ctx,
		Close 		: cancel,
	}
	go pool.loop()
	return pool
}

func (pool *Pool) loop() {

	interval := time.Duration(conf.GetSys().ClearInterval) * time.Second
	glog.Infoln("loop interval : ", interval)
	L: for {

		select {
		case <-pool.ctx.Done():
			pool.clean()
			break L
		case <-time.After(interval):		//NOTE: 每隔 interval 时间，执行一次清理
			pool.Upgrade()
		}
	}
	glog.Infoln("loop down")
}

func (pool *Pool) loopRead(ctx context.Context, conn *Conn) {

	L: for {

		msg, err := conn.ReadMessage()
		switch err {
		case nil:
			pool.readhub.Push(msg)
			pool.Update(conn)
		default:
			glog.Warningln("conn read error : ", err, "...tokened : ", conn.GetToken())
			select {
			case <-ctx.Done():		//NOTE: 如果已经在外部gg了，不需要再调用Cancel，这个地方不要执行清理，清理放在Listen中执行
			default:				//NOTE: 如果没有done，则主动调用Cancel，结束Listen
				conn.Cancel()
			}
			break L
		}
	}
}


//Listen 监听conn
func (pool *Pool) Listen(conn *Conn) (err error) {

	var ctx context.Context
	ctx, conn.Cancel = context.WithCancel(pool.ctx)
	conn.HandleContext(ctx)				//NOTE: ConnHandler 处理ctx

	glog.Infoln("after HandleContext")

	go pool.loopRead(ctx, conn)

	//NOTE: 下面的代码两个作用，1 阻塞，2 处理关闭
	select {

	case <-ctx.Done():					//NOTE: Context 被关闭了，执行清理
		pool.Remove(conn)
		err = errors.New("conn removed")
	}
	return
}


//Post 推送Message
func (pool *Pool) Post(msg *pb.Message) {

	pool.writehub.Push(msg)
}

//clean 清理Pool
func (pool *Pool) clean() {

	pool.writehub = nil
	pool.readhub = nil
}