package pool

import (
	"context"
	"errors"

	"github.com/golang/glog"
)

//connHub 用于存储管理conn
type connHub struct {

	req 		chan *req					//NOTE: 操作请求
	hub 		map[string]*conn			//NOTE: map[tokenString]*conn
}

const (
	typeAdd 	= iota	//NOTE: 增
	typeDel				//NOTE: 删
	typeSend			//NOTE: 操作[改?]
	typeFind			//NOTE: 查
)

type req struct {

	tag 	int					//NOTE: 请求类型
	param 	interface{}			//NOTE: 请求参数
	res 	chan interface{}	//NOTE: 结果chan
}

func (r *req) response(rlt interface{}) {

	if r.res != nil {
		r.res <- rlt
	}
}

func (r *req) close() {

	if r.res != nil {
		close(r.res)
	}
}

//newConnHub 构建ConnHub
func newConnHub(ctx context.Context) *connHub {

	ch := &connHub{
		req : make(chan *req),
		hub : make(map[string]*conn),
	}
	go ch.loop(ctx)
	return ch
}


func (ch *connHub) loop(ctx context.Context) {

	L: for {

		select {
		case <-ctx.Done(): break L
		case req := <-ch.req:
			switch req.tag {
			case typeAdd: ch.add(req)
			case typeDel: ch.del(req)
			case typeFind: ch.find(req)
			case typeSend: ch.send(req)
			default: 
				glog.Warning("request cann't knowing")
				req.close()			//NOTE: 避免奇怪的请求发过来，无法响应
			}
		}  
	}
}

func (ch *connHub) add(r *req) {

	c := r.param.(*conn)
	ch.hub[c.ToString()] = c
	r.response(true)
}

func (ch *connHub) del(r *req) {

	delete(ch.hub, r.param.(string))
	r.response(true)
}

func (ch *connHub) find(r *req) {

	r.response(ch.hub[r.param.(string)])
}


func (ch *connHub) send(r *req) {

	msg := r.param.(*Message)
	if c, ok := ch.hub[msg.GetToken()]; ok {
		go r.response(c.Send(msg))
		return
	}
	r.response(errors.New("cann't find conn tokened " + msg.GetToken()))
}
