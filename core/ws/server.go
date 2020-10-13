package ws

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/hiank/think/core"
	"github.com/hiank/think/settings"
)

//Server websocket server
type Server struct {
	upgrader   *websocket.Upgrader //NOTE: use default options
	*core.Pool                     //NOTE: 连接池
}

func newServer(ctx context.Context, msgHandler core.MessageHandler) *Server {

	return &Server{
		upgrader: new(websocket.Upgrader),
		Pool:     core.NewPool(ctx), //context.WithValue(ctx, core.CtxKeyRecvHandler, msgHandler)),
	}
}

//ServeHTTP 处理http 服务
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	tokenArr := r.Header["Token"]
	if len(tokenArr) == 0 {
		http.Error(w, "Non token component of the query", http.StatusNonAuthoritativeInfo) //NOTE: 没有包含token
		return
	}
	tokenStr := tokenArr[0]
	if !s.auth(tokenStr) {
		http.Error(w, "token auth fataled", http.StatusUnauthorized) //NOTE: token 认证失败
		return
	}

	defer core.Recover(core.Warning)
	c, err := s.upgrader.Upgrade(w, r, nil)
	core.Panic(err)
	defer c.Close()

	// core.Panic(s.Listen(token.GetBuilder().Get(tokenStr), &Handler{Conn: c, tokenStr: tokenStr}))
}

//auth 认证token
func (s *Server) auth(tokenStr string) bool {

	return true
}

var _singleServer *Server //NOTE: 全局唯一的websocket server

//Writer 写消息对象
type Writer int

//Handle 实现pool.MessageHandler
func (w Writer) Handle(msg core.Message) error {

	if _singleServer != nil {
		glog.Fatalf("websocket server not started, please start a websocket server first. (use 'ListenAndServe' function to do this.)")
	}
	return <-_singleServer.Push(msg)
}

//ListenAndServe used to start websocket serve NOTE: 只会有一个http服务被启动
//msgHandler 用于处理收到的消息
func ListenAndServe(ctx context.Context, ip string, msgHandler core.MessageHandler) (err error) {

	if _singleServer != nil {
		err = errors.New("websocket server existed, cann't start another one")
		glog.Fatal(err)
		return
	}

	_singleServer = newServer(ctx, msgHandler)
	// defer _singleServer.Close()

	http.Handle("/ws", _singleServer)
	server := &http.Server{Addr: core.WithPort(ip, settings.GetSys().WsPort)}
	go func() {
		<-ctx.Done()
		server.Close()
	}()
	return server.ListenAndServe()
}
