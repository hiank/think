package rpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	snet "net"
	"strconv"
	"sync"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/adapter/rpc/pp"
	"github.com/hiank/think/net/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"k8s.io/klog/v2"
)

func defaultListenOptions() listenOptions {
	return listenOptions{
		addr: ":10250",
		rest: new(pp.UnimplementedPipeServer),
	}
}

type funcCloser func() error

func (gc funcCloser) Close() error {
	return gc()
}

type listener struct {
	pp.UnimplementedPipeServer
	rest   IREST
	linkPP chan net.IConn
	io.Closer
}

//NewListener new a rpc listener
//NOTE: default addr is ":10250"
func NewListener(ctx context.Context, opts ...ListenOption) net.IListener {
	dopts := defaultListenOptions()
	for _, opt := range opts {
		opt.apply(&dopts)
	}
	lis, err := new(snet.ListenConfig).Listen(ctx, "tcp", dopts.addr)
	if err != nil {
		panic(fmt.Errorf("cannot listen in %x: %x", dopts.addr, err))
	}
	srv, linkPP, once := grpc.NewServer(), make(chan net.IConn), new(sync.Once)
	l := &listener{
		rest:   dopts.rest,
		linkPP: linkPP,
		Closer: funcCloser(func() (err error) {
			once.Do(func() {
				close(linkPP)
				srv.Stop()
				err = lis.Close()
			})
			return
		}),
	}
	go func() {
		defer l.Close()
		pp.RegisterPipeServer(srv, l)
		srv.Serve(lis)
	}()
	return l
}

func (l *listener) Accept() (c net.IConn, err error) {
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

func (l *listener) Link(ls pp.Pipe_LinkServer) (err error) {
	if identity, suc := l.linkAuth(ls); suc {
		ctx, cancel := context.WithCancel(ls.Context())
		l.linkPP <- &conn{identity: identity, ctx: ctx, cancel: cancel, sr: ls}
		<-ctx.Done()
	}
	return errors.New("link closed")
}

func (l *listener) Get(ctx context.Context, req *pb.Carrier) (res *pb.Carrier, err error) {
	return l.rest.Get(ctx, req)
}

func (l *listener) Post(ctx context.Context, req *pb.Carrier) (res *emptypb.Empty, err error) {
	return l.rest.Post(ctx, req)
}
