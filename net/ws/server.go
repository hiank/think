package ws

import (
	"github.com/hiank/think/conf"
	"strconv"
	"bytes"
	"context"
	"github.com/golang/glog"
	"net/http"
	"github.com/gorilla/websocket"
	// "github.com/hiank/think/util"
	"github.com/hiank/think/pool"
)

var upgrader = websocket.Upgrader{}		//NOTE: use default options
func wsServer(w http.ResponseWriter, r *http.Request) {

	tokenArr := r.Header["Token"]
	if  len(tokenArr) == 0 {
		http.Error(w, "Non token component of the query", http.StatusNonAuthoritativeInfo)		//NOTE: 没有包含token
		return
	}
	token := tokenArr[0]
	glog.Infoln("ws token : ", token)
	// token := "1000"

	//NOTE: 验证token

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Warning("upgrade error : ", err)		//NOTE: 一般不会出现这个地方的错误
		return
	}
	defer c.Close()

	conn, err := pool.NewConn("ws", token, &handler{conn:c})
	if err != nil {
		glog.Warning("ws create conn : ", err)
		return
	}
	defer conn.GetToken().Cancel()

	GetWSPool().Push(conn)
	GetWSPool().Listen(conn)
}


//ListenAndServeWS used to start websocket serve NOTE: 只会有一个http服务被启动
func ListenAndServeWS(ctx context.Context, addr string, mh pool.MessageHandler) (err error) {

	wsCtx, cancel := context.WithCancel(ctx)
	defer cancel()		//NOTE: 关闭context
	
	InitWSPool(wsCtx, mh)
	
	var buffer bytes.Buffer
	buffer.WriteString(addr)
	buffer.WriteByte(':')
	buffer.WriteString(strconv.FormatInt(conf.GetSys().WsPort, 10))
	
	addr = buffer.String()

	http.HandleFunc("/ws", wsServer)
	if err = http.ListenAndServe(addr, nil); err != nil {

		glog.Error("listen websocket error " + err.Error())
	}
	return
}
