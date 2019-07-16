package k8s

import (
	"github.com/hiank/think/token"
	"github.com/golang/glog"
	"context"
	"github.com/hiank/think/pb"
	"github.com/hiank/think/pool"
	"sync"

	tg "github.com/hiank/think/net/k8s/protobuf"
)


var k8sconnpool *pool.Pool
var k8sCtx context.Context
var k8smu sync.Mutex
//InitK8SPool 初始化k8s conn pool
func InitK8SPool(ctx context.Context, h pool.MessageHandler) {

	k8smu.Lock()

	if k8sconnpool == nil {

		k8sCtx = ctx
		k8sconnpool = pool.NewPool(ctx, h)
	}

	k8smu.Unlock()
}

// GetK8SPool 获得全局k8sconnpool
func GetK8SPool() (cm *pool.Pool) {

	k8smu.Lock()

	if k8sconnpool != nil {

		select {
		case <-k8sCtx.Done():
			k8sconnpool = nil
		default:
			cm = k8sconnpool
		}
	}

	k8smu.Unlock()
	return
}

// CloseK8SPool clean the static object
func CloseK8SPool() {

	k8smu.Lock()

	if k8sconnpool != nil {

		select {
		case <-k8sCtx.Done():
		default:
			k8sconnpool.Close()
			k8sconnpool = nil	
		}
	}

	k8smu.Unlock()
}


//**********************************************//
type connhandler struct {

	conn 	tg.Pipe_LinkServer
}

func newConnHandler(lc tg.Pipe_LinkServer) pool.ConnHandler {

	si := &connhandler {
		conn 		: lc,
	}
	return si
}



func (si *connhandler) ReadMessage(ctx context.Context) (msg *pb.Message, err error) {

	if msg, err = si.conn.Recv(); err == nil {
		msg.Key = ctx.Value(token.ContextKey("key")).(string)
	}
	glog.Infoln("ReadMessage in connHandler key : ", msg)
	return
}


func (si *connhandler) WriteMessage(ctx context.Context, msg *pb.Message) (err error) {
 
	return si.conn.Send(msg)
}
