package pool

import (
	"github.com/golang/glog"
	"errors"
	"github.com/hiank/think/pb"
	"github.com/hiank/think/conf"
	"time"
	"context"
)

//Pool 消息处理的核心
type Pool struct {

	*ConnHub

	readhub 	*MessageHub 			//NOTE: 读消息hub
	writehub 	*MessageHub 			//NOTE: 写消息hub

	rb 			*runtineHub 			//NOTE: 用于存储, key 绑定的Context, 比如某个grpc服务出了问题,可以通过这个方法释放所有Conn资源

	ctx			context.Context	 		//NOTE: Context 用于维护生命周期
	Close 		context.CancelFunc 		//NOTE: 关闭方法
}


//NewPool 创建一个新的Pool
func NewPool(ctx context.Context, mh MessageHandler) *Pool {

	ctx, cancel := context.WithCancel(ctx)

	num := conf.GetSys().MaxMessageGo
	ch := newConnHub()
	pool := &Pool {
		readhub 	: NewMessageHub(mh, num),
		writehub 	: NewMessageHub(ch, num),
		rb 			: newRuntineHub(),
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

		msg, err := conn.Recv()
		select {
		case <-ctx.Done(): 				break L
		case <-conn.GetToken().Done(): 	break L
		default:
			switch err {
			case nil:
				pool.readhub.Push(msg)
				pool.Update(conn)
			default:
				glog.Warningln("conn read error : ", err, "...tokened : ", conn.GetToken().ToString())
				conn.GetToken().Cancel()
				break L
			}
		}
	}
}


//Listen 监听conn
func (pool *Pool) Listen(conn *Conn) (err error) {

	r := pool.rb.get(pool.ctx, conn.GetKey())
	go pool.loopRead(r.Context, conn)

	//NOTE: 下面的代码两个作用，1 阻塞，2 处理关闭
	select {

	case <-conn.GetToken().Done():		//NOTE: 此Conn 关联的Token 被释放了
		err = errors.New("conn tokened : " + conn.GetToken().ToString() + " Done")
	case <-r.Done():					//NOTE: Context 被关闭了，执行清理
		err = errors.New("conn keyed : " + conn.GetKey() + " Done")
	}
	pool.Remove(conn)
	return
}


//Post 推送Message
func (pool *Pool) Post(msg *pb.Message) {

	pool.writehub.Push(msg)
}

// //Cancel 释放key 为key 的Conn
// func (pool *Pool) Cancel(key string) {

// 	pool.rb.delete(key)
// }

//clean 清理Pool
func (pool *Pool) clean() {

	pool.writehub = nil
	pool.readhub = nil
}