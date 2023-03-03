package net

import (
	"context"
	"io"
	"sync"

	"github.com/hiank/think/run"
)

type Client struct {
	ctx    context.Context
	m      sync.Map
	router RouteMux
	dialer Dialer
	knower Knower
	io.Closer
}

func NewClient(ctx context.Context, dialer Dialer, knower Knower) *Client {
	cli := &Client{
		dialer: dialer,
		knower: knower,
	}
	cli.ctx, cli.Closer = run.StartHealthyMonitoring(ctx)
	return cli
}

// AutoSend 自动发送消息
// NOTE: 执行的步骤为
//  1. 消息自动解析出服务名
//  2. 从缓存中找到client
//     a. 未找到client，创建并初始化一个->缓存到clientset中
//  3. 使用client的Send发送此消息
//
// NOTE: 如果client的第一个任务(连接到服务器)失败，会执行client的Close，导致此client被从clientset中移除. 期间缓存的消息将全部丢失
func (cli *Client) AutoSend(tm *Message) (err error) {
	addr, err := cli.knower.ServeAddr(tm)
	if err == nil {
		err = cli.autoClientConnset(addr).Send(tm)
	}
	return
}

// RouteMux for register Handler to handle message recved
func (cli *Client) RouteMux() *RouteMux {
	return &cli.router
}

func (cli *Client) autoClientConnset(addr string) (ccp *clientConnset) {
	v, loaded := cli.m.LoadOrStore(addr, &clientConnset{})
	if ccp = v.(*clientConnset); !loaded {
		ccp.set = newConnset(cli.RouteMux())
		ccp.dial = func(ctx context.Context) (Conn, error) {
			return cli.dialer.Dial(ctx, addr)
		}
		ccp.ctx, ccp.Closer = run.StartHealthyMonitoring(cli.ctx, ccp.set.close, func() {
			cli.m.Delete(addr)
		})
	}
	return
}

// clientConnset 客户连接集
type clientConnset struct {
	ctx  context.Context
	dial func(context.Context) (Conn, error)
	set  *connset
	io.Closer
}

func (ccp *clientConnset) Send(tm *Message) (err error) {
	lc, err := ccp.set.loadOrStore(ccp.ctx, tm.Token().ToString(), func(ctx context.Context) (tc Conn, err error) {
		if tc, err = ccp.dial(ctx); err != nil {
			ccp.Close() //dial not work. remove the client from clientset
		}
		return
	})
	if err == nil {
		if err = lc.Send(tm); err != nil {
			lc.Close() //connect failed or disconnected
		}
	}
	return
}
