package rpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	snet "net"
	"strconv"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/adapter/rpc/pp"
	"github.com/hiank/think/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"k8s.io/klog/v2"
)

func defaultListenOptions() listenOptions {
	return listenOptions{
		addr: ":10250",
		rest: new(pp.UnimplementedPipeServer),
	}
}

type listener struct {
	pp.UnsafePipeServer //for PipeServcer
	REST                //for PipeServer
	io.Closer
	linkPP chan net.IAC
}

//NewListener new a rpc listener
//NOTE: default addr is ":10250"
func NewListener(ctx context.Context, opts ...ListenOption) net.Listener {
	dopts := defaultListenOptions()
	for _, opt := range opts {
		opt.apply(&dopts)
	}
	lis, err := new(snet.ListenConfig).Listen(ctx, "tcp", dopts.addr)
	if err != nil {
		panic(fmt.Errorf("cannot listen in %x: %x", dopts.addr, err))
	}
	srv, linkPP := grpc.NewServer(), make(chan net.IAC)
	l := &listener{
		REST:   dopts.rest,
		linkPP: linkPP,
		Closer: run.NewOnceCloser(func() error {
			close(linkPP)
			srv.Stop()
			return lis.Close()
		}),
	}
	go func() {
		defer l.Close()
		pp.RegisterPipeServer(srv, l)
		srv.Serve(lis)
	}()
	return l
}

func (l *listener) Accept() (c net.IAC, err error) {
	c, ok := <-l.linkPP
	if !ok {
		err = io.EOF
	}
	return
}

//linkAuth get 'identity' value from metadata
//NOTE: identity is generated in grpc client. it was generated with key "hostname.uid" in redis (IStorage)
func (l *listener) linkAuth(ls pp.Pipe_LinkServer) (identity uint64, suc bool) {
	if md, ok := metadata.FromIncomingContext(ls.Context()); ok {
		if arr := md.Get("identity"); arr != nil || len(arr) > 0 {
			if identity, err := strconv.ParseUint(arr[0], 10, 64); err == nil {
				ls.SetHeader(metadata.Pairs("success", "true"))
				if err = ls.SendHeader(nil); err == nil {
					return identity, true
				}
			}
		}
	}
	klog.Warning("identity metadata for grpc invalid")
	return
}

//Link for pp.PipeServer
func (l *listener) Link(ls pp.Pipe_LinkServer) (err error) {
	if identity, suc := l.linkAuth(ls); suc {
		ctx, cancel := context.WithCancel(ls.Context())
		l.linkPP <- net.IAC{ID: strconv.FormatUint(identity, 10), Conn: &conn{ctx: ctx, cancel: cancel, s: ls}}
		<-ctx.Done()
	}
	return errors.New("link closed")
}
