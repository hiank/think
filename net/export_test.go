package net

import "github.com/hiank/think/net/pb"

var (
	Export_newTestConn = func(recvPP <-chan pb.M, sendPP chan<- pb.M) Conn {
		return &testConn{
			recvPP: recvPP,
			sendPP: sendPP,
		}
	}
)
