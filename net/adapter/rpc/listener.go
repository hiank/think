package rpc

import (
	"context"
	"fmt"
	"io"
	snet "net"
	"strconv"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/adapter"
	"github.com/hiank/think/net/adapter/rpc/pipe"
	"github.com/hiank/think/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	ErrLinkClosed = run.Err("rpc: Link closed")

	linkMetadataIdentity = "identity"
	linkMetadataSuccess  = "success"
)

type listener struct {
	addr      string //listen address with port
	keepalive pipe.KeepaliveServer
	rest      pipe.RestServer
	// healthy   *run.Healthy
	pipe.UnsafeKeepaliveServer
	io.Closer
	adapter.ChanAccepter
}

//Link for LinkServer
func (lis *listener) Link(ls pipe.Keepalive_LinkServer) (err error) {
	identity, err := linkAuth2(ls)
	if err == nil {
		ctx, cancel := context.WithCancel(ls.Context())
		defer cancel()
		lis.ChanAccepter <- net.IdentityConn{ID: strconv.FormatUint(identity, 10), Conn: &conn{ctx: ctx, cancel: cancel, s: ls}}
		<-ctx.Done()
		err = ErrLinkClosed
	}
	return
}

//servePipe serve pipeServer
func (lis *listener) servePipe(slis snet.Listener) {
	defer lis.Close()
	srv := grpc.NewServer()
	defer srv.Stop()
	pipe.RegisterKeepaliveServer(srv, lis.keepalive)
	pipe.RegisterRestServer(srv, lis.rest)
	srv.Serve(slis)
}

//linkAuth get 'identity' value from metadata
//NOTE: identity is generated in grpc client. it was generated with key "hostname.uid" in redis (IStorage)
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

type ListenOption run.Option[*listener]

func defaultListener() listener {
	return listener{
		keepalive: pipe.UnimplementedKeepaliveServer{},
		rest:      pipe.UnimplementedRestServer{},
		addr:      ":30202",
		// healthy:      run.NewHealthy(),
		ChanAccepter: make(adapter.ChanAccepter),
	}
}

func WithAddress(addr string) ListenOption {
	return run.FuncOption[*listener](func(lis *listener) {
		lis.addr = addr
	})
}

func WithDefaultKeepaliveServer() ListenOption {
	return run.FuncOption[*listener](func(lis *listener) {
		lis.keepalive = lis
	})
}

func WithRestServer(rest pipe.RestServer) ListenOption {
	return run.FuncOption[*listener](func(lis *listener) {
		lis.rest = rest
	})
}

func NewListener(ctx context.Context, opts ...ListenOption) net.Listener {
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
	go lis.servePipe(slis)
	return &lis
}
