package ws

import (
	"github.com/gorilla/websocket"
)

var (
	Export_newConn = func(wc *websocket.Conn) *conn {
		return &conn{
			wc: wc,
		}
	}
	Export_withDefaultListenOption = withDefaultListenOption
	Export_newUnimplementedAuther  = func() *unimplementedAuther {
		return &unimplementedAuther{}
	}
)
