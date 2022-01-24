package rpc

type listenOptions struct {
	addr string
	rest REST
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

func WithREST(rest REST) ListenOption {
	return funcListenOption(func(lo *listenOptions) {
		lo.rest = rest
	})
}
