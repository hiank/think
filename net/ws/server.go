package ws

import (
	"context"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/hiank/think/pb"
	"github.com/hiank/think/pool"
	"github.com/hiank/think/settings"
	"github.com/hiank/think/token"
	"github.com/hiank/think/utils"
	"github.com/hiank/think/utils/robust"
)

//Server websocket server
type Server struct {

	upgrader *websocket.Upgrader	//NOTE: use default options
	*pool.Pool						//NOTE: 连接池
}

func newServer(ctx context.Context, msgHandler pool.MessageHandler) *Server {

	return &Server {
		upgrader 	: new(websocket.Upgrader),
		Pool 		: pool.NewPool(context.WithValue(ctx, pool.CtxKeyRecvHandler, msgHandler)),
	}
}

//ServeHTTP 处理http 服务
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	tokenArr := r.Header["Token"]
	if  len(tokenArr) == 0 {
		http.Error(w, "Non token component of the query", http.StatusNonAuthoritativeInfo)		//NOTE: 没有包含token
		return
	}
	tokenStr := tokenArr[0]
	if !s.auth(tokenStr) {
		http.Error(w, "token auth fataled", http.StatusUnauthorized)		//NOTE: token 认证失败
		return
	}

	defer robust.Recover(robust.Warning)
	c, err := s.upgrader.Upgrade(w, r, nil)
	robust.Panic(err)
	defer c.Close()

	tok, err := token.GetBuilder().Get(tokenStr)
	robust.Panic(err)

	robust.Panic(s.AddConn(pool.NewConn(tok, &Handler{Conn:c, tokenStr: tokenStr})))
}

//auth 认证token
func (s *Server) auth(tokenStr string) bool {

	return true
}


var _singleServer *Server		//NOTE: 全局唯一的websocket server

//Writer 写消息对象
type Writer int 

//Handle 实现pool.MessageHandler
func (w Writer) Handle(msg *pb.Message) error {

	if _singleServer != nil {
		glog.Fatalf("websocket server not started, please start a websocket server first. (use 'ListenAndServe' function to do this.)")
	}
	_singleServer.Post(msg)
	return nil
}

//ListenAndServe used to start websocket serve NOTE: 只会有一个http服务被启动
//msgHandler 用于处理收到的消息
func ListenAndServe(ctx context.Context, ip string, msgHandler pool.MessageHandler) (err error) {

	if _singleServer != nil {
		glog.Fatal("websocket server existed, cann't start another one.")
	}

	_singleServer = newServer(ctx, msgHandler)
	defer _singleServer.Close()

	http.Handle("/ws", _singleServer)
	return http.ListenAndServe(utils.WithPort(ip, settings.GetSys().WsPort), nil)
}
