package rpc

import (
	"context"
	"fmt"
	"io"
	snet "net"
	"strconv"

	"github.com/hiank/think/auth"
	"github.com/hiank/think/net"
	"github.com/hiank/think/net/adapter"
	"github.com/hiank/think/net/adapter/rpc/pipe"
	"github.com/hiank/think/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"k8s.io/klog/v2"
)

const (
	ErrLinkClosed = run.Err("rpc: Link closed")

	linkMetadataIdentity = "identity"
	linkMetadataSuccess  = "success"
)

type funcLinkServer func(tkey string) (auth.Token, chan<- net.Conn)

func (fls funcLinkServer) Link(ls pipe.Keepalive_LinkServer) error {
	identity, err := linkAuth2(ls)
	if err != nil {
		return err
	}
	tk, cc := fls(strconv.FormatUint(identity, 10))
	ctx, closer := run.StartHealthyMonitoring(ls.Context(), func() {
		tk.Close()
	})
	cc <- &conn{
		tk:     tk,
		sr:     ls,
		Closer: closer,
	}
	<-ctx.Done()
	return ErrLinkClosed
}

type listenOptions struct {
	addr      string
	rest      pipe.RestServer
	keepalive pipe.KeepaliveServer
	accepter  adapter.ChanAccepter
}

type ListenOption run.Option[*listenOptions]

func defaultListenOptions() listenOptions {
	return listenOptions{
		keepalive: new(pipe.UnimplementedKeepaliveServer),
		rest:      new(pipe.UnimplementedRestServer),
		addr:      ":30202",
		accepter:  make(adapter.ChanAccepter),
	}
}

func WithAddress(addr string) ListenOption {
	return run.FuncOption[*listenOptions](func(lis *listenOptions) {
		lis.addr = addr
	})
}

func WithServeKeepalive(ts auth.Tokenset) ListenOption {
	return run.FuncOption[*listenOptions](func(opts *listenOptions) {
		opts.keepalive = funcLinkServer(func(tkey string) (auth.Token, chan<- net.Conn) {
			return ts.Derive(tkey), opts.accepter
		})
	})
}

func WithRestServer(rest pipe.RestServer) ListenOption {
	return run.FuncOption[*listenOptions](func(opts *listenOptions) {
		opts.rest = rest
	})
}

type listener struct {
	io.Closer
	adapter.ChanAccepter
}

// linkAuth get 'identity' value from metadata
// NOTE: identity is generated in grpc client. it was generated with key "hostname.uid" in redis (IStorage)
func linkAuth2(ls pipe.Keepalive_LinkServer) (identity uint64, err error) {
	err = fmt.Errorf("rpc: identity metadata for grpc-link invalided")
	if md, ok := metadata.FromIncomingContext(ls.Context()); ok {
		if arr := md.Get(linkMetadataIdentity); arr != nil || len(arr) > 0 {
			if identity, err = strconv.ParseUint(arr[0], 10, 64); err == nil {
				if err = ls.SetHeader(metadata.Pairs(linkMetadataSuccess, "true")); err == nil {
					err = ls.SendHeader(nil)
				}
			}
		}
	}
	return
}

// NewListener for grpc
func NewListener(ctx context.Context, opts ...ListenOption) net.Listener {
	dopts := defaultListenOptions()
	for _, opt := range opts {
		opt.Apply(&dopts)
	}
	slis, err := new(snet.ListenConfig).Listen(ctx, "tcp", dopts.addr)
	if err != nil {
		panic(err) //failed listen in given address
	}
	lis := &listener{ChanAccepter: dopts.accepter}
	_, lis.Closer = run.StartHealthyMonitoring(ctx, func() {
		close(dopts.accepter)
		slis.Close()
	})
	go func() {
		defer lis.Close()
		srv := grpc.NewServer()
		defer srv.Stop()
		pipe.RegisterKeepaliveServer(srv, dopts.keepalive)
		pipe.RegisterRestServer(srv, dopts.rest)
		klog.Warning(srv.Serve(slis))
	}()
	return lis
}
