package k8s


import (
	"github.com/hiank/think/token"
	"os"
	"github.com/golang/glog"
	"google.golang.org/grpc/connectivity"
	"errors"
	"google.golang.org/grpc"
	"github.com/hiank/think/pb"
	"context"
	"github.com/hiank/think/pool"
	"sync"

	tg "github.com/hiank/think/net/k8s/protobuf"
)


var clientpool *pool.Pool
var poolCtx context.Context
var cpmu sync.RWMutex

//InitClientPool 初始化clientpool
func InitClientPool(ctx context.Context, h pool.MessageHandler) {

	cpmu.Lock()

	if clientpool == nil {

		poolCtx = ctx
		clientpool = pool.NewPool(ctx, h)
	}

	cpmu.Unlock()
}


// GetClientPool get initialized static value
func GetClientPool() (cm *pool.Pool) {

	cpmu.RLock()

	if clientpool != nil {

		select {
		case <-poolCtx.Done():
			clientpool = nil
		default:
			cm = clientpool
		}
	}

	cpmu.RUnlock()
	return
}

// CloseClientPool clean the static object 
func CloseClientPool() {

	cpmu.Lock()

	if clientpool != nil {

		select {
		case <-poolCtx.Done():
		default:
			clientpool.Close()
		}
		clientpool = nil
	}

	cpmu.Unlock()
}




//*****************************************************//
//ClientHandler type
const (
	TypeNormal 			= iota 		//NOTE: 一次性调用
	TypeStream 						//NOTE: 流
)

//ClientHandler grpc 客户端读写
type ClientHandler struct {
	// pool.Identifier

	// ctx 		context.Context			//NOTE: 

	linkPool 	*pool.Pool				//NOTE:
	rChan 		chan *pb.Message		//NOTE: 

	cc 			*grpc.ClientConn		//NOTE: gprc 客户端连接
	client		tg.PipeClient			//NOTE: 定义的grpc 客户端
}

//NewClientHandler 创建一个新的ClientHandler
func NewClientHandler(cc *grpc.ClientConn) *ClientHandler {

	return &ClientHandler{
		rChan 		: make(chan *pb.Message),
		cc 			: cc,
		client 		: tg.NewPipeClient(cc),
	}	
}

//WriteMessage 向grpc 远端发送Message
func (ch *ClientHandler) WriteMessage(ctx context.Context, msg *pb.Message) (err error) {

	var t int
	if t, err = pb.GetServerType(msg.GetData()); err != nil {
		glog.Warningln(err)
		return
	}

	switch t {
	case pb.TypeGET:
		var res *pb.Message
		if res, err = ch.client.Get(ctx, msg); err == nil {
			ch.rChan <- res
		}
	case pb.TypePOST:
		_, err = ch.client.Post(ctx, msg)
	case pb.TypeSTREAM:
		if ch.linkPool == nil {
			ch.linkPool = pool.NewPool(ctx, &linkReadHandler{ch.rChan})
			go ch.checkHealth(ctx)
		}
		if !ch.linkPool.CheckConnected(ctx.Value(token.ContextKey("key")).(string), ctx.Value(token.ContextKey("token")).(string)) {
			errChan := make(chan error)
			go ch.listenLink(ctx, msg, errChan)
			var ok bool
			if err, ok = <-errChan; ok {
				glog.Warning(err)
				break
			}
		}
		ch.linkPool.Post(msg)
	default: err = errors.New("cann't operate message type undefined")
	}
	return
}

func (ch *ClientHandler) listenLink(ctx context.Context, msg *pb.Message, errChan chan error) {

	lp := ch.linkPool
	lc, err := ch.client.Link(ctx)
	if err != nil {
		errChan <- err
		return
	}
	defer lc.CloseSend()				//NOTE: 退出时关闭link

	hostname := os.Getenv("HOSTNAME")
	lc.Send(&pb.Message{Key: hostname, Token: msg.GetToken()})

	conn := pool.NewConnWithDerivedToken(msg.GetKey(), msg.GetToken(), &linkClientHandler{conn:lc})
	defer conn.GetToken().Cancel()		//NOTE: 退出时执行清理
	lp.Push(conn)
	close(errChan)						//NOTE: 如果一切正常，关闭errChan
	lp.Listen(conn)
}

// //stream 处理TypeSTREAM 数据 写
// func (ch *ClientHandler) streamWrite(ctx context.Context, msg *pb.Message) (err error) {

// 	lp := ch.linkPool
// 	if !lp.CheckConnected(ctx.Value(token.ContextKey("key")).(string), ctx.Value(token.ContextKey("token")).(string)) {

// 		var lc tg.Pipe_LinkClient
// 		if lc, err = ch.client.Link(ctx); err != nil {
// 			glog.Warningln(err)
// 			return
// 		}

// 		hostname := os.Getenv("HOSTNAME")
// 		lc.Send(&pb.Message{Key: hostname, Token: msg.GetToken()})
// 		lh := &linkClientHandler {
// 			// Identifier 	: pool.NewDefaultIdentifier(msg.GetKey(), msg.GetToken()),
// 			conn 		: lc,
// 		}
// 		// conn := pool.NewDefaultConn(lh)
// 		conn := pool.NewConnWithDerivedToken(msg.GetKey(), msg.GetToken(), lh)
// 		lp.Push(conn)

// 		wait := make(chan bool)
// 		go func() {
// 			close(wait)
// 			lp.Listen(conn)
// 		} ()
// 		<-wait
// 	}
// 	lp.Post(msg)
// 	return
// }

//ReadMessage 从grpc 远端读取数据
//如果返回一个错误，则Pool 将感知到这个Conn 出了问题，会做相应处理
func (ch *ClientHandler) ReadMessage(ctx context.Context) (msg *pb.Message, err error) {

	var ok bool 
	if msg, ok = <- ch.rChan; !ok {

		err = errors.New("k8s client read chan closed")
	}
	return
}

//checkHealth 健康检查，注意 cc 肯定经历过Ready 的状态，才可能逻辑上执行到这一步
func (ch *ClientHandler) checkHealth(ctx context.Context) {

	L: for {
		s := ch.cc.GetState()
		switch s {
		case connectivity.Shutdown: fallthrough
		case connectivity.TransientFailure:
			break L
		}
		if !ch.cc.WaitForStateChange(ctx, s) {
			break L
		}
	}
	//NOTE: 此处表明 cc 连接出现问题
	ch.linkPool.Close()		//NOTE: 关闭连接池
	close(ch.rChan)			//NOTE: 关闭读chan，产生一个读错误，触发Pool的读错误处理
}


//*****************************************************//
type linkReadHandler struct {

	rChan 	chan *pb.Message		//NOTE: 读到的数据传入到这个chan中
}

//Handle 处理从grpc conn中读到的Message
func (lh linkReadHandler) Handle(m *pb.Message) error {

	defer func() {
		if r := recover(); r != nil {
			glog.Warning(r)
		}
	}()
	lh.rChan <- m					//NOTE: 这个chan 可能会被外部close，用于关闭conn的读，因为消息处理是异步的，可能存在之前读到的消息延后处理的情况，导致向关闭下chan写消息的情况，要处理pannic
	return nil
}


//*****************************************************//

type linkClientHandler struct {

	conn 	tg.Pipe_LinkClient
}

//WriteMessage 向连接中写入数据
func (lh *linkClientHandler) WriteMessage(ctx context.Context, msg *pb.Message) (err error) {

	return lh.conn.Send(msg)
}

//ReadMessage 从连接中读取数据
func (lh *linkClientHandler) ReadMessage(ctx context.Context) (msg *pb.Message, err error) {

	if msg, err = lh.conn.Recv(); err == nil {
		// msg.Key = lh.GetKey()
		msg.Key = ctx.Value(token.ContextKey("key")).(string)
	}
	return
}

