package net

var (
	Export_newTestConn = func(recvPP <-chan *Doc, sendPP chan<- *Doc) Conn {
		return &testConn{
			recvPP: recvPP,
			sendPP: sendPP,
		}
	}
)
