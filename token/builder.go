package token

import (
	"github.com/golang/glog"
	"errors"
	"sync"
	"context"
)

type tokenBuilder struct {

	ctx 	context.Context				//NOTE: Builder 的基础Context
	Cancel 	context.CancelFunc

	hub 	map[string]*Token 			//NOTE: map[tokenStr]*Token
}


func (tb *tokenBuilder) get(tokenStr string) (*Token, bool) {

	token, ok := tb.hub[tokenStr]
	if !ok {

		glog.Infoln("tokenBuilder cann't get token : ", tokenStr)
		for key := range tb.hub {
			glog.Infoln("tokenBuilder get : ", key)
		}
	}
	return token, ok
}

func (tb *tokenBuilder) build(tokenStr string) (token *Token, err error) {

	if _, ok := tb.hub[tokenStr]; ok {
		err = errors.New("token '" + tokenStr + "' existed in cluster")
		return
	}
	token, _ = newToken(context.WithValue(tb.ctx, ContextKey("token"), tokenStr))		//NOTE: 此处一定不会触发error
	tb.hub[tokenStr] = token
	glog.Infoln("tokenBuilder build token : ", tokenStr)
	for key := range tb.hub {
		glog.Infoln("tokenBuilder build : ", key)
	}
	return
}

func (tb *tokenBuilder) delete(tokenStr string) {

	glog.Infoln("tokenBuilder delete token : ", tokenStr)
	delete(tb.hub, tokenStr)
}


var builder *tokenBuilder
var mtx sync.Mutex

//InitBuilder 获取单例的tokenBuilder
func InitBuilder(ctx context.Context) {

	mtx.Lock()
	if builder == nil {

		builder = &tokenBuilder{
			hub : make(map[string]*Token),
		}
		builder.ctx, builder.Cancel = context.WithCancel(ctx)
	}
	mtx.Unlock()
}

//Get 根据字符token 找到Token，如果不存在，则新建一个
func Get(tokenStr string) (token *Token, ok bool, err error) {

	mtx.Lock()
	defer mtx.Unlock()

	if builder == nil {
		err = errors.New("package token error : without initialized builder. please call InitBulider(context.Context) first")
		return
	}

	select {
	case <-builder.ctx.Done():
		err = builder.ctx.Err()
		builder.Cancel()
		builder = nil
	default:
		token, ok = builder.get(tokenStr)
	}
	return
}

//Build 创建并返回一个Token
func Build(tokenStr string) (token *Token, err error) {

	mtx.Lock()
	defer mtx.Unlock()

	if builder == nil {
		err = errors.New("package token error : without initialized builder. please call InitBulider(context.Context) first")
		return
	}

	select {
	case <-builder.ctx.Done():
		err = builder.ctx.Err()
		builder.Cancel()
		builder = nil
	default:
		token, err = builder.build(tokenStr)
	}
	return
}

//CloseBuilder 清除
func CloseBuilder() {

	mtx.Lock()
	if builder != nil {
		builder.Cancel()
		builder = nil
	}
	mtx.Unlock()
}
