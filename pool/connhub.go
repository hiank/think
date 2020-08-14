package pool

import (
	"context"
	"errors"

	"github.com/hiank/think/settings"
)

//connHub 用于存储管理conn
type connHub struct {
	hub     map[string]*Conn //NOTE: map[tokenString]*Conn
	getReq  chan *connGetReq
	sendReq chan *connSendReq
	rmReq   chan string
}

//connGetReq '获取'请求
type connGetReq struct {
	tokStr  string
	rlt     chan<- *Conn
	builder connBuilder
}

//connSendReq 发送消息请求
type connSendReq struct {
	msg *Message
	err chan<- error
}

//newConnHub 构建ConnHub
func newConnHub(ctx context.Context) *connHub {

	ch := &connHub{
		getReq:  make(chan *connGetReq, settings.GetSys().ConnHubReqLen),
		sendReq: make(chan *connSendReq, settings.GetSys().ConnHubReqLen),
		rmReq:   make(chan string, settings.GetSys().ConnHubReqLen),
		hub:     make(map[string]*Conn),
	}
	go ch.loop(ctx)
	return ch
}

func (ch *connHub) loop(ctx context.Context) {

L:
	for {
		select {
		case <-ctx.Done():
			break L
		case tokStr := <-ch.rmReq:
			delete(ch.hub, tokStr)
		case req := <-ch.getReq:
			ch.responseGetReq(req)
		case req := <-ch.sendReq:
			ch.responseSendReq(req)
		}
	}
}

//responseGetReq 响应获取请求
//如果hub中没有请求的conn，则根据是否包含connBuilder 来确定是否要构建一个conn
//查找和构建 都使用这个方法来处理
func (ch *connHub) responseGetReq(req *connGetReq) {

	c, ok := ch.hub[req.tokStr]
	if !ok {
		if req.builder == nil {
			close(req.rlt)
			return
		}
		c = req.builder()
		ch.hub[req.tokStr] = c
	}
	req.rlt <- c
}

func (ch *connHub) responseSendReq(req *connSendReq) {

	handleError := func(err error) {
		if req.err != nil {
			req.err <- err
		}
	}
	c, ok := ch.hub[req.msg.GetToken()]
	if !ok {
		handleError(errors.New("cann't find Conn tokened " + req.msg.GetToken()))
		return
	}
	go handleError(c.Send(req.msg))
}
