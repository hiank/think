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
	ctx      context.Context
	Close    context.CancelFunc
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

//build 构建Conn
func (p *Pool) build(tok *tk.Token, rw IO) (*Conn, error) {

	select {
	case <-tok.Done():
		return nil, errors.New("token was invalided")
	default:
	}

	rlt := make(chan *Conn)
	p.getReq <- &connGetReq{
		tokStr: tok.ToString(),
		rlt:    rlt,
		builder: connBuilder(func() *Conn {
			return newConn(tok, rw)
		}),
	}
	return <-rlt, nil
}

//Listen 启动监听
func (p *Pool) Listen(tok *tk.Token, rw IO) error {

	conn, err := p.build(tok, rw)
	if err != nil {
		return err
	}
	err = conn.Listen(p.ctx.Value(CtxKeyRecvHandler).(MessageHandler))
	p.rmReq <- conn.ToString() //NOTE: 监听结束后，删除引用
	return err
}

//PostAndWait 推送消息，等待反馈
func (p *Pool) PostAndWait(msg *Message) error {

	err := make(chan error)
	p.sendReq <- &connSendReq{msg, err}
	return <-err
}

//Post 推送消息
func (p *Pool) Post(msg *Message) {

	p.sendReq <- &connSendReq{msg, nil}
}
