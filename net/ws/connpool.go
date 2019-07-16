package ws

import (
	"github.com/hiank/think/token"
	"github.com/golang/glog"
	"context"
	"github.com/hiank/think/pb"
	"github.com/gorilla/websocket"
	"github.com/hiank/think/pool"
	"sync"
)

var wsconnpool *pool.Pool
var poolCtx context.Context
var cpmu sync.RWMutex

// InitWSPool 初始化ws pool
func InitWSPool(ctx context.Context, rh pool.MessageHandler) {

	cpmu.Lock()

	if wsconnpool == nil {

		poolCtx = ctx
		wsconnpool = pool.NewPool(ctx, rh)
	}
	
	cpmu.Unlock()
}

// GetWSPool get initialized static value
func GetWSPool() (cm *pool.Pool) {

	cpmu.RLock()

	if wsconnpool != nil {

		select {
		case <-poolCtx.Done():		//NOTE: 如果context 已经关闭，表明pool 会被关闭，需要处理这个情况
			wsconnpool = nil
		default:
			cm = wsconnpool
		}
	}

	cpmu.RUnlock()
	return
}

// CloseWSPool clean the static object 
func CloseWSPool() {

	cpmu.Lock()

	if wsconnpool != nil {

		select {
		case <-poolCtx.Done():		//NOTE: 如果已经外部关闭了，不需要再调用pool 的Close
		default:
			wsconnpool.Close()
		}
		wsconnpool = nil
	}

	cpmu.Unlock()
}


//**********************************************//

type handler struct {

	// pool.Identifier
	// pool.IgnoreHandleContext

	conn 	*websocket.Conn
}


//ReadMessage 读消息，实现frame.Conn
func (c *handler) ReadMessage(ctx context.Context) (msg *pb.Message, err error) {

	_, buf, err := c.conn.ReadMessage()		//NOTE: 从websocket 读取消息
	if err == nil {	
		glog.Infoln("ws conn read message :", buf)
		m, err := pb.AnyDecode(buf)
		if err == nil {			//NOTE: 解析消息

			glog.Infoln("ws conn any decode :", m)
			if key, err := pb.GetServerKey(m); err == nil {
				msg = &pb.Message{Key: key, Token: ctx.Value(token.ContextKey("token")).(string), Data: m}
			}
		}
	}
	return
}


// WriteMessage Writer
func (c *handler) WriteMessage(ctx context.Context, msg *pb.Message) (err error) {


	var buf []byte
	if buf, err = pb.AnyEncode(msg.GetData()); err != nil {
		glog.Infoln(err)
		return
	}
	err = c.conn.WriteMessage(websocket.BinaryMessage, buf)
	switch (err) {
	case nil:
	default:			//NOTE:	处理错误
		glog.Infoln(err)
	}
	return
}

