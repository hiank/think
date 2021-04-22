package token

import "time"

type options struct {
	timeout   time.Duration
	needCache bool
}

type Option interface {
	apply(*options)
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

//WithTimeout 设置超时
func WithTimeout(val time.Duration) Option {
	return newFuncOption(func(opts *options) {
		opts.timeout = val
	})
}

func WithCache() Option {
	return newFuncOption(func(opts *options) {
		opts.needCache = true
	})
}
