package net

import "google.golang.org/protobuf/types/known/anypb"

var (
	Export_newTestConn = func(identity uint64, recvPP <-chan *anypb.Any, sendPP chan<- *anypb.Any) Conn {
		return &testConn{
			identity: identity,
			recvPP:   recvPP,
			sendPP:   sendPP,
		}
	}
)
