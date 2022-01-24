package ws

import (
	"github.com/gorilla/websocket"
	"github.com/hiank/think/net"
)

var (
	Export_NewConn = func(uid uint64, wc *websocket.Conn) net.Conn {
		return &conn{uid: uid, wc: wc}
	}
)
