package ws

import (
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hiank/think/net"
	"github.com/hiank/think/oauth"
	"k8s.io/klog/v2"
)

//ServeHelper websocket连接核心
type ServeHelper struct {
	server   *http.Server
	upgrader *websocket.Upgrader //NOTE: use default options
	auther   oauth.Auther
	dopts    *options
	net.ChanAccepter
}

//NewServeHelper 新建一个ServeHelper
func NewServeHelper(addr string, auther oauth.Auther, opts ...HelperOption) *ServeHelper {
	dopts := newDefaultOptions()
	for _, opt := range opts {
		opt.apply(dopts)
	}
	helper := &ServeHelper{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		auther:       auther,
		dopts:        dopts,
		ChanAccepter: make(net.ChanAccepter),
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
	uid, err := helper.auther.Auth(tokenArr[0])
	if err != nil {
		http.Error(w, "token auth fataled:"+err.Error(), http.StatusUnauthorized) //NOTE: token 认证失败
		return
	}

	wsc, err := helper.upgrader.Upgrade(w, r, nil)
	switch err {
	case nil:
		var c net.Conn
		if helper.dopts.connMaker != nil { //NOTE: 如果有自定义的ConnMaker，使用其生成Conn
			c = helper.dopts.connMaker.Make()
		} else { //NOTE: 否则生成默认的Conn
			c = &conn{ReadWriteCloser: wsc, uid: uid}
		}
		helper.ChanAccepter <- c
	case io.EOF:
	default:
		klog.Warning(err)
	}
}

//ListenAndServe 启动服务
func (helper *ServeHelper) ListenAndServe() error {
	return helper.server.ListenAndServe()
}

//Close 关闭
func (helper *ServeHelper) Close() error {
	close(helper.ChanAccepter)
	return helper.server.Close()
}
