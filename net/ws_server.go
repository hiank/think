package net

import (
	"github.com/hiank/think/util"
	"io"
	// "fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/golang/glog"
)


var upgrader = websocket.Upgrader{} // use default options
func wsServer(w http.ResponseWriter, r *http.Request) {


	token := r.FormValue("token")
	if token == "" {
		http.Error(w, "Non token component of the query", http.StatusNonAuthoritativeInfo)	//NOTE: 没有包含token
		return
	}

	//NOTE: 验证token

	defer util.RecoverErr("wsServer : ")

	c, err := upgrader.Upgrade(w, r, nil)
	util.PanicErr(err)

	quit := make(chan bool)
	conn := NewConn(c, token)
	// GetConnPool().GetConnaddChan() <- conn	//NOTE: 将此连接发送给connpool，添加
	GetConnPool().AddConn(conn) 			//NOTE: 将conn添加到ConnPool中

	go conn.SendAsync(quit)

	k8schan := GetK8sClient().GetK8sRequestChan()
L:	for {
		select {
		case <-quit: break L
		default:
			msg, err := conn.ReadMessage()				//NOTE: 读取并处理客户端发来的数据
			GetConnPool().UpdateConn(conn)				//NOTE: 更新conn在ConnPool中的连接状态
			switch err {
			case nil: k8schan <- msg
			case io.EOF: break L
			}
		}
	}

	close(quit)
	// GetConnPool().GetConndelChan() <- conn	//NOTE: 将此连接发送到connpool，删除
	GetConnPool().DelConn(conn)				//NOTE: 将conn从ConnPool中删除
}


//ListenAndServeWS used to start websocket serve
func ListenAndServeWS(addrWithPort string) (err error) {

	connpool, client := GetConnPool(), GetK8sClient()
	defer func() {

		ReleaseConnPool()
		ReleaseK8sClient()
	}()

	quit := make(chan bool)
	defer close(quit)

	go connpool.RecvAsync(quit, client.GetK8sResponseChan())	//NOTE: 接受集群返回的消息，处理

	http.HandleFunc("/websocket", wsServer)
	err = http.ListenAndServe(addrWithPort, nil)
	if err != nil {
		glog.Error("listen websocket error " + err.Error())
	}
	return
}
