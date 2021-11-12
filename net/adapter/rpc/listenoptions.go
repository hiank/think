package rpc

type listenOptions struct {
	addr string
	rest IREST
}

type ListenOption interface {
	apply(*listenOptions)
}

type funcListenOption func(*listenOptions)

func (flo funcListenOption) apply(lo *listenOptions) {
	flo(lo)
}

func WithAddr(addr string) ListenOption {
	return funcListenOption(func(lo *listenOptions) {
		lo.addr = addr
	})
}

func WithREST(rest IREST) ListenOption {
	return funcListenOption(func(lo *listenOptions) {
		lo.rest = rest
	})
}
