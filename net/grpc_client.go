package net

import (
	"github.com/hiank/think/cluster"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	tg "github.com/hiank/think/net/protobuf/grpc"

	"github.com/hiank/think/util"
	proto "github.com/hiank/think/net/protobuf"
)


// K8sRequest 接受grpc调用的请求信息
type K8sRequest interface {

	GetToken() (token []byte)
	GetData() (data *any.Any)
}

// K8sResponse 接受grpc调用返回到结果信息
type K8sResponse interface {

	GetTokens() (tokens [][]byte)
	GetData() (data *any.Any)
}


type k8sConn struct {

	name 	string 				//NOTE: 连接的服务名
	addr 	string 				//NOTE: 连接的具体服务器地址

	conn 	*grpc.ClientConn
	pipe 	tg.Pipe_TranClient
	tokens 	map[string]byte		//Note: tokens store token which connecting to the server
}

// newK8sConn create a new k8sConn object and initialized
func newK8sConn(serverName string, addr string, conn *grpc.ClientConn) (kc *k8sConn) {

	defer util.RecoverErr("newK8sConn error : ")

	grpclient := tg.NewPipeClient(conn)
	pipe, err := grpclient.Tran(context.Background())
	util.PanicErr(err)

	kc = &k8sConn{

		name 	: serverName,
		addr 	: addr,
		conn	: conn,
		pipe 	: pipe,
		tokens	: make(map[string]byte),
	}
	return
}

func (kc *k8sConn) getServerName() string {

	return kc.name
}

func (kc *k8sConn) getServerAddr() string {

	return kc.addr
}

func (kc *k8sConn) include(token string) bool {

	_, ok := kc.tokens[token]
	return ok
}

func (kc *k8sConn) bind(token string) {

	kc.tokens[token] = '0'
}

func (kc *k8sConn) unbind(token string) {

	delete(kc.tokens, token)
}

func (kc *k8sConn) recv(ch chan K8sResponse, exit chan *k8sConn) {

	defer func() {
		util.RecoverErr("RecvSync error : ")
		exit <- kc
	}()

	pipe := kc.pipe

	for {

		msg, err := pipe.Recv()
		util.PanicErr(err)

		ch <- msg
	}
}

// close 当K8sClient 删除对此连接的管理时，调用
func (kc *k8sConn) close() {

	kc.pipe.CloseSend()
	kc.conn.Close()
}



//**********************************K8sClient under************************************


// K8sClient manager k8s cluster's communication
type K8sClient struct {

	inchan 		chan K8sRequest		//NOTE: 用于接收消息，转发到k8s中
	upchan 		chan K8sResponse	//NOTE: 用于转发k8s返回的消息
	close 		chan *k8sConn
	quit 		chan bool			//NOTE: 数据处理协程退出chan

	k8sconns 	map[string]map[string]*k8sConn	// Note: map[protoName/serverName]map[ip]*k8sConn
}

var k8sclient *K8sClient
// GetK8sClient get initialized static value
func GetK8sClient() (kc *K8sClient) {

	if k8sclient != nil {

		kc = k8sclient
		return
	}

	kc = &K8sClient{
		inchan		: make(chan K8sRequest),
		upchan		: make(chan K8sResponse),
		close 		: make(chan *k8sConn),
		quit 		: make(chan bool),
		k8sconns	: make(map[string]map[string]*k8sConn),
	}

	go kc.async()
	k8sclient = kc
	return
}


// ReleaseK8sClient release and clear static client object
func ReleaseK8sClient() {

	if k8sclient == nil {
		return
	}

	close(k8sclient.inchan)
	close(k8sclient.upchan)
	close(k8sclient.close)
	close(k8sclient.quit)

	for _, v := range k8sclient.k8sconns {

		for _, k8sconn := range v {

			k8sconn.close()
		}
	}
	k8sclient.k8sconns = nil
	k8sclient = nil
}


func (kc *K8sClient) async() {

L:	for {
		select {
		case <-kc.quit: return
		case req := <-kc.inchan:	//NOTE: 向远端发送消息协程
			anyMsg := req.GetData()
			name, err := ptypes.AnyMessageName(anyMsg)
		
			grpconn, err := kc.dial(req.GetToken(), name)
			if err != nil {
				continue L
			}
			grpconn.pipe.Send(req.(*tg.Request))
		case k8sconn := <-kc.close:	//NOTE: 当k8sconn协程退出的时候，维护列表中删除此数据
			kc.removeConn(k8sconn)
			k8sconn.close()
		}
	}
}


// GetK8sRequestChan 向服务发送消息，通过此chan
func (kc *K8sClient) GetK8sRequestChan() chan K8sRequest {

	return kc.inchan
}

// GetK8sResponseChan 服务返回消息后，将向此chan写入
func (kc *K8sClient) GetK8sResponseChan() chan K8sResponse {

	return kc.upchan
}


// dial make grpc connection between client and server
func (kc *K8sClient) dial(token []byte, name string) (*k8sConn, error) {

	defer util.RecoverErr("dail grpc : ")

	name = proto.GetServerName(name)
	tokenStr := string(token)
	k8sconns, ok := kc.k8sconns[name]
	if ok {

		for _, conn := range k8sconns {

			if conn.include(tokenStr) {

				return conn, nil
			}
		}
	} else {

		k8sconns = make(map[string]*k8sConn)
		kc.k8sconns[name] = k8sconns
	}


	addr, err := cluster.GetAddr(cluster.TypeKubIn, name)
	util.PanicErr(err)

	k8sconn, ok := k8sconns[addr]
	if !ok {

		grpconn, err := grpc.Dial(addr, grpc.WithInsecure())
		util.PanicErr(err)

		k8sconn = newK8sConn(name, addr, grpconn)
		go k8sconn.recv(kc.upchan, kc.close)	//NOTE: 启动监听协程，接收远端发送的消息，推送到upchan
		k8sconns[addr] = k8sconn
	}
	k8sconn.bind(tokenStr)
	return k8sconn, nil
}


func (kc *K8sClient) removeConn(k8sconn *k8sConn) {

	m, ok := kc.k8sconns[k8sconn.getServerName()]
	if !ok {
		return
	}

	_, ok = m[k8sconn.getServerAddr()]
	if !ok {
		return
	}

	delete(m, k8sconn.getServerAddr())
	if len(m) == 0 {
		delete(kc.k8sconns, k8sconn.getServerName())
	}

}

