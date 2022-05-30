package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/adapter/rpc/pipe"
	"github.com/hiank/think/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	ErrLinkAuthFailed = run.Err("rpc: link auth failed")
)

type DialOption run.Option[*keepaliveDialer]

func WithIdentity(id string) DialOption {
	return run.FuncOption[*keepaliveDialer](func(kd *keepaliveDialer) {
		kd.identity = id
	})
}

func WithTimeout(timeout time.Duration) DialOption {
	return run.FuncOption[*keepaliveDialer](func(kd *keepaliveDialer) {
		kd.timeout = timeout
	})
}

type keepaliveDialer struct {
	identity string
	timeout  time.Duration
}

func NewKeepaliveDialer(opts ...DialOption) net.Dialer {
	kd := &keepaliveDialer{
		timeout: time.Second * 10,
	}
	for _, opt := range opts {
		opt.Apply(kd)
	}
	return kd
}

//keep-alive connect dial
func (kd *keepaliveDialer) Dial(ctx context.Context, addr string) (c net.Conn, err error) {
	dialCtx, dialCancel := context.WithTimeout(ctx, time.Second*10)
	defer dialCancel()
	cc, err := grpc.DialContext(dialCtx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name))) //NOTE: block 为阻塞直到ready，insecure 为不需要验证的
	if err != nil {
		return
	}
	lc, err := pipe.NewKeepaliveClient(cc).Link(metadata.NewOutgoingContext(ctx, metadata.Pairs(linkMetadataIdentity, kd.identity)))
	if err == nil {
		var md metadata.MD
		if md, err = lc.Header(); err == nil {
			if ss := md.Get(linkMetadataSuccess); len(ss) > 0 && ss[0] == "true" {
				ctx, closer := run.StartHealthyMonitoring(ctx, run.CloserToDoneHook(cc))
				c = &conn{
					ctx:    ctx,
					s:      lc,
					Closer: closer,
				}
				return
			}
			err = ErrLinkAuthFailed
		}
	}
	return
}

func RestDial(ctx context.Context, addr string) (cli RestClient, err error) {
	dialCtx, dialCancel := context.WithTimeout(ctx, time.Second*10)
	defer dialCancel()
	cc, err := grpc.DialContext(dialCtx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name))) //NOTE: block 为阻塞直到ready，insecure 为不需要验证的
	if err == nil {
		rc := &restClient{RestClient: pipe.NewRestClient(cc)}
		_, rc.Closer = run.StartHealthyMonitoring(ctx, run.CloserToDoneHook(cc))
		cli = rc
	}
	return
}
