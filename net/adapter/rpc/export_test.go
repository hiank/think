package rpc

import (
	"context"
	snet "net"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/adapter/rpc/pipe"
	"github.com/hiank/think/run"
	"google.golang.org/grpc"
)

var (
	Export_defaultListener  = defaultListener
	Export_convertolistener = func(lis net.Listener) *listener {
		return lis.(*listener)
	}
	Export_NewListenerEx = func(ctx context.Context, opts ...ListenOption) (*listener, snet.Listener, *grpc.Server) {
		lis := defaultListener()
		for _, opt := range opts {
			opt.Apply(&lis)
		}
		ctx, cancel := context.WithCancel(ctx)
		slis, err := new(snet.ListenConfig).Listen(ctx, "tcp", lis.addr)
		if err != nil {
			panic(err) //failed listen in given address
		}
		healthy := run.NewHealthy()
		lis.Closer = run.NewHealthyCloser(healthy, cancel)
		go healthy.Monitoring(ctx, func() {
			close(lis.ChanAccepter)
			slis.Close()
		}) //monitor ctx
		srv := grpc.NewServer()
		go func() {
			defer lis.Close()
			defer srv.Stop()
			pipe.RegisterKeepaliveServer(srv, lis.keepalive)
			pipe.RegisterRestServer(srv, lis.rest)
			srv.Serve(slis)
		}()
		return &lis, slis, srv
	}
)
