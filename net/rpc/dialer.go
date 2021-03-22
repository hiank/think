package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/hiank/think/net"
	tg "github.com/hiank/think/net/rpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

//Dialer grpc连接器
//10秒钟超时
var Dialer = dialerFunc(func(ctx context.Context, target string) (conn net.Conn, err error) {

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	cc, err := grpc.DialContext(ctxWithTimeout, target, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name))) //NOTE: block 为阻塞直到ready，insecure 为不需要验证的
	if err == nil {
		conn = newPipe(ctx, target, tg.NewPipeClient(cc))
	}
	return
})

type dialerFunc func(ctx context.Context, target string) (conn net.Conn, err error)

func (df dialerFunc) Dial(ctx context.Context, target string) (conn net.Conn, err error) {
	return df(ctx, target)
}
