package adapter

import "context"

type ListenerOption interface {
	apply(*listenerOptions)
}

type listenerOptions struct {
	ctx context.Context
	// ip   string
	// port int
	addr string
}

type funcListenerOption func(*listenerOptions)

func (fl funcListenerOption) apply(los *listenerOptions) {
	fl(los)
}

// func WithIP(ip string) ListenerOption {
// 	return funcListenerOption(func(lo *listenerOptions) {
// 		lo.ip = ip
// 	})
// }

// func WithPort(port int) ListenerOption {
// 	return funcListenerOption(func(lo *listenerOptions) {
// 		lo.port = port
// 	})
// }

func WithAddr(addr string) ListenerOption {
	return funcListenerOption(func(lo *listenerOptions) {
		lo.addr = addr
	})
}

func WithContext(ctx context.Context) ListenerOption {
	return funcListenerOption(func(lo *listenerOptions) {
		lo.ctx = ctx
	})
}

// func WithAddr(addr string) Lis
