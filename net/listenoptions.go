package net

import (
	"context"

	"github.com/hiank/think/pool"
)

type listenOptions struct {
	ctx         context.Context
	connHandler func(Conn)   //NOTE: accept得到的conn的处理方法
	recvHandler pool.Handler //NOTE: 默认connHandler中收到的消息的处理方法
}

func newDefaultListenOptions() *listenOptions {
	return &listenOptions{
		ctx: context.Background(),
	}
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcListenOption struct {
	f func(*listenOptions)
}

func (fo *funcListenOption) apply(opts *listenOptions) {
	fo.f(opts)
}

func newFuncOption(f func(*listenOptions)) *funcListenOption {
	return &funcListenOption{
		f: f,
	}
}

// ListenOption configures how we set up the connection.
type ListenOption interface {
	apply(*listenOptions)
}

//WithConnHandler 自定义的Conn处理方法
//默认的Conn处理，是将Conn放到pool中集中管理
func WithConnHandler(handler func(Conn)) ListenOption {
	return newFuncOption(func(opts *listenOptions) {
		opts.connHandler = handler
	})
}

//WithContext 自定义的Context
//默认的Context，为context.Background()
func WithContext(ctx context.Context) ListenOption {
	return newFuncOption(func(opts *listenOptions) {
		opts.ctx = ctx
	})
}

//WithRecvHandler 收到的消息的处理方法
//如果未自定义ConnHandler则务必设置此选项
//如果自定义了ConnHandler则此选项无效
func WithRecvHandler(handler pool.Handler) ListenOption {
	return newFuncOption(func(opts *listenOptions) {
		opts.recvHandler = handler
	})
}
