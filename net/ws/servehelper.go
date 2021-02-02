package ws

import (
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hiank/think/net"
	"k8s.io/klog/v2"
)

//ServeHelper websocket连接核心
type ServeHelper struct {
	server   *http.Server
	upgrader *websocket.Upgrader //NOTE: use default options
	connChan chan net.Conn
	net.Accepter
}

//NewServeHelper 新建一个ServeHelper
func NewServeHelper(addr string) *ServeHelper {

	ch := make(chan net.Conn, 8)
	helper := &ServeHelper{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}, //new(websocket.Upgrader),
		connChan: ch,
		Accepter: net.ChanAccepter(ch),
	}

	httpHandler := http.NewServeMux()
	httpHandler.Handle("/ws", helper)
	helper.server = &http.Server{Addr: addr, Handler: httpHandler}
	return helper
}

//ServeHTTP 处理http 服务
func (helper *ServeHelper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	tokenArr := r.Header["Token"]
	if len(tokenArr) == 0 {
		http.Error(w, "Non token component of the query", http.StatusNonAuthoritativeInfo) //NOTE: 没有包含token
		return
	}
	uid, ok := helper.auth(tokenArr[0])
	if !ok {
		http.Error(w, "token auth fataled", http.StatusUnauthorized) //NOTE: token 认证失败
		return
	}

	c, err := helper.upgrader.Upgrade(w, r, nil)
	switch err {
	case nil:
		helper.connChan <- &conn{ReadWriteCloser: c, uid: uid}
	case io.EOF:
	default:
		klog.Warning(err)
	}
}

//auth 认证token
func (helper *ServeHelper) auth(tokenStr string) (uint64, bool) {

	return 1001, true
}

//ListenAndServe 启动服务
func (helper *ServeHelper) ListenAndServe() error {

	return helper.server.ListenAndServe()
}

//Close 关闭
func (helper *ServeHelper) Close() error {

	close(helper.connChan)
	return helper.server.Close()
}
