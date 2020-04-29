package net

import (
	"context"

	"github.com/hiank/think/net/k8s"
	"github.com/hiank/think/net/ws"
	"github.com/hiank/think/pool"
)

//ServeK8s 启动一个k8s服务，同一个进程只能有一个k8s服务
func ServeK8s(ip string, msgHandler k8s.MessageHandler) error {

	return k8s.ListenAndServe(GetRuntime().Context, ip, msgHandler)
}

//ServeWS 启动一个ws服务，同一个进程只能有一个ws服务
//收到的消息交给k8s ClientHub 来处理
func ServeWS(ip string) error {

	return ws.ListenAndServe(GetRuntime().Context, ip, k8s.NewClientHub(context.WithValue(GetRuntime().Context, pool.CtxKeyRecvHandler, new(ws.Writer))))
}
