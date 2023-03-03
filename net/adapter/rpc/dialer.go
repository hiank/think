package rpc

import (
	"context"
	"fmt"
	"io"

	"github.com/hiank/think/auth"
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

// type DialOption run.Option[*keepaliveDialer]

// func WithTimeout(timeout time.Duration) DialOption {
// 	return run.FuncOption[*keepaliveDialer](func(kd *keepaliveDialer) {
// 		kd.timeout = timeout
// 	})
// }

type keepaliveDialer struct {
	tk auth.Token
	// timeout time.Duration
}

func NewKeepaliveDialer(tk auth.Token) net.Dialer {
	return &keepaliveDialer{tk}
}

// keep-alive connect dial
// @param ctx: only for dial to remote server. after connected, ignore it's 'Done()'
func (kd *keepaliveDialer) Dial(ctx context.Context, addr string) (c net.Conn, err error) {
	// dialCtx, dialCancel := context.WithTimeout(ctx, time.Second*10)
	// defer dialCancel()
	cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name))) //NOTE: block 为阻塞直到ready，insecure 为不需要验证的
	if err != nil {
		return
	}
	lc, err := pipe.NewKeepaliveClient(cc).Link(metadata.NewOutgoingContext(ctx, metadata.Pairs(linkMetadataIdentity, kd.tk.ToString())))
	if err != nil {
		return
	}
	md, err := lc.Header()
	if err != nil {
		return
	}
	if ss := md.Get(linkMetadataSuccess); len(ss) > 0 && ss[0] == "true" {
		_, closer := run.StartHealthyMonitoring(kd.tk, run.CloserToDoneHook(cc))
		c = &conn{
			tk:     kd.tk,
			sr:     lc,
			Closer: closer,
		}
		return
	}
	err = ErrLinkAuthFailed
	return
}

type RestClient interface {
	pipe.RestClient
	io.Closer
}

type restClient struct {
	pipe.RestClient
	io.Closer
}

// RestDial for new RestClient
func RestDial(ctx context.Context, addr string) (cli RestClient, err error) {
	// dialCtx, dialCancel := context.WithTimeout(ctx, time.Second*10)
	// defer dialCancel()
	cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name))) //NOTE: block 为阻塞直到ready，insecure 为不需要验证的
	if err == nil {
		rc := restClient{RestClient: pipe.NewRestClient(cc)}
		_, rc.Closer = run.StartHealthyMonitoring(ctx, run.CloserToDoneHook(cc))
		cli = rc
	}
	return
}
