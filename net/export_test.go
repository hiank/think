package net

import (
	"context"

	"github.com/hiank/think/pool"
)

var (
	Export_newDefaultListenOptions = newDefaultListenOptions
	Export_ListenOptionApply       = func(opt ListenOption) func(*listenOptions) {
		return opt.apply
	}
	Export_getListenOptionsCtx = func(opts *listenOptions) context.Context {
		return opts.ctx
	}
	Export_getListenOptionsConnHandler = func(opts *listenOptions) func(Conn) {
		return opts.connHandler
	}
	Export_getListenOptionsRecvHandler = func(opts *listenOptions) pool.Handler {
		return opts.recvHandler
	}
	Export_getClientCtx = func(client *Client) context.Context {
		return client.ctx
	}
	Export_getClientDialer = func(client *Client) Dialer {
		return client.dialer
	}
	Export_getClientRecvHandler = func(client *Client) pool.Handler {
		return client.recvHandler
	}
	Export_getClientHubPool = func(client *Client) *pool.HubPool {
		return client.hubPool
	}
	Export_ClientAutoHub = func(client *Client) func(key string) *pool.Hub {
		return client.autoHub
	}
	Export_LoopRecv = loopRecv
	// Export_LoopRecv = func(client *Client) func(ctx context.Context, reciver Reciver, handler pool.Handler) {
	// 	return loopRecv
	// }
)
