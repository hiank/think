package token

import (
	"context"
	"sync"

	"github.com/hiank/think/settings"
)

//tokenReq 构建请求
type tokenReq struct {
	tokenStr  string
	tok       chan<- *Token
	withBuild bool //NOTE: 如果没有找到时，是否需要构建
}

//Builder 用于构建Token
type Builder struct {
	ctx   context.Context   //NOTE: Builder 的基础Context
	hub   map[string]*Token //NOTE: map[tokenStr]*Token
	rmReq chan *Token       //NOTE: 已Done 的token
	req   chan *tokenReq    //NOTE: 请求构建token
}

func newBuilder() *Builder {

	builder := &Builder{
		ctx:   BackgroundLife().Context,
		hub:   make(map[string]*Token),
		rmReq: make(chan *Token, settings.GetSys().TokDonedLen), //NOTE: 设置缓存，避免清理token 阻塞
		req:   make(chan *tokenReq, settings.GetSys().TokReqLen),
	}
	go builder.healthMonitoring()
	return builder
}

//healthMonitoring 监测状态
func (b *Builder) healthMonitoring() {

	for {
		select {
		case tok := <-b.rmReq: //NOTE: 已Done 的token
			b.delete(tok)
		case req := <-b.req:
			b.response(req)
		}
	}
}

//removeReq 删除请求
func (b *Builder) removeReq() chan<- *Token {

	return b.rmReq
}

//Get get *Token and if cann't find the *Token, then Build one and return
func (b *Builder) Get(tokenStr string) *Token {

	tokRes := make(chan *Token)
	b.req <- &tokenReq{
		tokenStr:  tokenStr,
		tok:       tokRes,
		withBuild: true,
	}
	return <-tokRes
}

//Find 查找Token
func (b *Builder) Find(tokenStr string) (tok *Token, ok bool) {

	tokRes := make(chan *Token)
	b.req <- &tokenReq{
		tokenStr:  tokenStr,
		tok:       tokRes,
		withBuild: false,
	}
	tok, ok = <-tokRes
	return
}

//response 处理请求，注意这个是在专门的goroutine 中执行的，不存在数据竞争问题
func (b *Builder) response(req *tokenReq) {

	tok, ok := b.hub[req.tokenStr]
	if !ok {
		if !req.withBuild {
			close(req.tok)
			return
		}
		tok, _ = newToken(context.WithValue(b.ctx, IdentityKey, req.tokenStr)) //NOTE: 此处一定不会触发error
		b.hub[req.tokenStr] = tok
	}
	req.tok <- tok
}

//delete 删除传入的token
func (b *Builder) delete(tok *Token) {

	if tok == nil {
		return
	}
	str := tok.ToString()
	// if b.hub[tok.ToString()] == tok {
	delete(b.hub, str)
	// }
}

var _singleBuilder *Builder
var _singleBuilderOnce sync.Once

//GetBuilder return singleton tokenBuilder object
func GetBuilder() *Builder {

	_singleBuilderOnce.Do(func() {

		_singleBuilder = newBuilder()
	})
	return _singleBuilder
}
