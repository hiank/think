package pool

import (
	"context"
	"errors"

	tk "github.com/hiank/think/token"
)

//contextKey type of context key
type contextKey int

//CtxKeyRecvHandler 用于处理收到的消息，当连接池中的连接收到消息后，传入这个key对应的chan，有pool来处理
var CtxKeyRecvHandler = contextKey(0)

//Pool 连接池
type Pool struct {
	*connHub //NOTE: 处理建立的连接

	ctx   context.Context
	Close context.CancelFunc
}

//NewPool 构建Pool
//ctx must contained 'CtxKeyRecvHandler'
func NewPool(ctx context.Context) *Pool {

	p := &Pool{
		connHub: newConnHub(ctx),
	}
	p.ctx, p.Close = context.WithCancel(ctx)
	return p
}

//Listen 监听conn
func (p *Pool) Listen(tok *tk.Token, rw IO, addedArr ...chan interface{}) (err error) {

	conn := newConn(tok, rw)
	select {
	case <-conn.Done():
		err = errors.New("conn Done")
		if len(addedArr) > 0 {
			addedArr[0] <- false
		}
	default:
		var added chan interface{}
		if len(addedArr) > 0 {
			added = addedArr[0]
		}
		p.req <- &req{tag: typeAdd, param: conn, res: added}
		err = conn.Listen(p.ctx.Value(CtxKeyRecvHandler).(MessageHandler))
	}
	return
}

//Has 查找Conn
func (p *Pool) Has(tokStr string) bool {

	req := &req{tag: typeFind, param: tokStr, res: make(chan interface{})}
	p.req <- req
	return (<-req.res).(bool)
}

//PostAndWait 推送消息，等待反馈
func (p *Pool) PostAndWait(msg *Message) error {

	req := &req{tag: typeSend, param: msg, res: make(chan interface{})}
	p.req <- req
	return (<-req.res).(error)
}

//Post 推送消息
func (p *Pool) Post(msg *Message) error {

	p.req <- &req{tag: typeSend, param: msg}
	return nil
}
