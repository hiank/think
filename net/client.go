package net

import (
	"context"

	"github.com/hiank/think/run"
)

const (
	ErrUnimplementApi = run.Err("net: unimplement api")
)

//Client 客户端，每种类型的协议使用唯一的客户端
type Client struct {
	dialer Dialer
	knower Knower
	route  *RouteMux
	pool   *connpool
}

//NewClient new Client. for auto dial to server
func NewClient(ctx context.Context, dialer Dialer, knower Knower) *Client {
	cli := &Client{
		dialer: dialer,
		knower: knower,
		route:  &RouteMux{},
	}
	cli.pool = newConnpool(ctx, cli.route)
	return cli
}

//RouteMux for register Handler
func (cli *Client) RouteMux() *RouteMux {
	return cli.route
}

//Send 自动处理消息
//NOTE: 执行的步骤为
//1. 消息自动解析出服务名
//2. 从缓存中找到Hub
//3. 未找到Hub
//4. 新建一个Hub，并尝试建立连接
//5. 连接失败，删除Hub
//implement Sender interface
func (cli *Client) AutoSend(im IdentityMessage) (err error) {
	id, err := cli.knower.Identity(im.ID)
	if err == nil {
		var tc *taskConn
		if tc, err = cli.pool.loadOrStore(id, func(ctx context.Context) (c Conn, err error) {
			addr, err := cli.knower.ServeAddr(im.M)
			if err == nil {
				c, err = cli.dialer.Dial(ctx, addr)
			}
			return
		}); err == nil {
			err = tc.Send(im.M)
		}
	}
	return
}

func (cli *Client) Close() error {
	return cli.pool.Close()
}
