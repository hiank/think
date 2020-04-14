package pool

import (
	"context"
	"errors"

	"github.com/hiank/think/pb"
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

	t 	int					//NOTE: 请求类型
	r 	interface{}			//NOTE: 请求参数
	s 	chan interface{}	//NOTE: 结果chan
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
		case req := <-ch.req: ch.do(req)
		}  
	}
}

func (ch *connHub) do(r *req) {

	var s interface{}
	switch r.t {
	case typeAdd:
		c := r.r.(*conn)
		ch.hub[c.ToString()] = c
		s = true
	case typeDel:
		delete(ch.hub, r.r.(string))
		s = true
	case typeSend:
		msg := r.r.(*pb.Message)
		if c, ok := ch.hub[msg.GetToken()]; ok {		
			if errChan := c.AsyncSend(r.r.(*pb.Message)); r.s != nil {
				go func(res chan interface{}, errChan <-chan error) {
					res <- (<- errChan)
				}(r.s, errChan)
			}
		} else {
			s = errors.New("cann't find conn tokened " + msg.GetToken())
		}
	case typeFind:
		s = ch.hub[r.r.(string)]
	}

	if (s != nil) && (r.s != nil) {
		r.s <- s
	}
}
