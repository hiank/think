package ws

import "github.com/hiank/think/net"

type options struct {
	connMaker ConnMaker
}

func newDefaultOptions() *options {
	return &options{}
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcOption struct {
	f func(*options)
}

func (fo *funcOption) apply(opts *options) {
	fo.f(opts)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// HelperOption configures how we set up the connection.
type HelperOption interface {
	apply(*options)
}

//WithConnMaker 自定义的Conn构造器
//默认的Conn构造方法是，生成*conn
func WithConnMaker(maker ConnMaker) HelperOption {
	return newFuncOption(func(opts *options) {
		opts.connMaker = maker
	})
}

type ConnMaker interface {
	Make() net.Conn
}
