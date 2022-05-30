package net

import (
	"context"
	"io"
	"sync"

	"github.com/hiank/think/net/box"
	"github.com/hiank/think/net/one"
	"github.com/hiank/think/run"
)

const (
	ErrUnimplementedApi = run.Err("net: unimplemented api")
)

//client 客户端，每种类型的协议使用唯一的客户端
type client struct {
	ctx  context.Context
	dial func(context.Context) (Conn, error)
	set  *connset
	io.Closer
}

func (cli *client) Send(tm TokenMessage) (err error) {
	id := tm.Token.Value(box.ContextkeyTokenUid).(string)
	lc, err := cli.set.loadOrStore(cli.ctx, id, func(ctx context.Context) (tc TokenConn, err error) {
		if tc.T, err = cli.dial(ctx); err == nil {
			tc.Token = one.TokenSet().Derive(id)
		} else {
			cli.Close() //dial not work. remove the client from clientset
		}
		return
	})
	if err == nil {
		if err = lc.Send(tm.T); err != nil {
			lc.Close() //connect failed or disconnected
		}
	}
	return
}

type clientset struct {
	ctx    context.Context
	m      sync.Map
	router RouteMux
	dialer Dialer
	knower Knower
	io.Closer
}

func NewClientset(ctx context.Context, dialer Dialer, knower Knower) Clientset {
	cs := &clientset{
		dialer: dialer,
		knower: knower,
	}
	cs.ctx, cs.Closer = run.StartHealthyMonitoring(ctx)
	return cs
}

//AutoSend 自动发送消息
//NOTE: 执行的步骤为
//1. 消息自动解析出服务名
//2. 从缓存中找到client
//	a. 未找到client，创建并初始化一个->缓存到clientset中
//3. 使用client的Send发送此消息
//NOTE: 如果client的第一个任务(连接到服务器)失败，会执行client的Close，导致此client被从clientset中移除. 期间缓存的消息将全部丢失
func (cs *clientset) AutoSend(tm TokenMessage) (err error) {
	addr, err := cs.knower.ServeAddr(tm.T)
	if err == nil {
		err = cs.autoClient(addr).Send(tm)
	}
	return
}

//RouteMux for register Handler to handle message recved
func (cs *clientset) RouteMux() *RouteMux {
	return &cs.router
}

func (cs *clientset) autoClient(addr string) (cli *client) {
	v, loaded := cs.m.LoadOrStore(addr, &client{})
	if cli = v.(*client); !loaded {
		cli.set = newConnset(cs.RouteMux())
		cli.dial = func(ctx context.Context) (Conn, error) {
			return cs.dialer.Dial(ctx, addr)
		}
		cli.ctx, cli.Closer = run.StartHealthyMonitoring(cs.ctx, cli.set.close, func() {
			cs.m.Delete(addr)
		})
	}
	return
}
