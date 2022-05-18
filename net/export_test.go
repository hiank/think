package net

import "github.com/hiank/think/net/box"

// import "github.com/hiank/think/net/pb"

var (
	Export_newTmpConn = func(recvPP chan *box.Message, sendPP chan<- *box.Message) Conn {
		return &tmpConn{
			recvPP: recvPP,
			sendPP: sendPP,
		}
	}
	// Export_newTaskConn = newTaskConn
)

// func Export_taskConnCtx(tc Conn) context.Context {
// 	return tc.(*taskConn).ctx
// }
