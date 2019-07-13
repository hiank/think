package net

import (
	"context"
	"github.com/hiank/think/net/k8s"
)

//ServeK8s 启动一个k8s服务
func ServeK8s(ctx context.Context, addr string, h k8s.MessageHandler) error {

	return k8s.ListenAndServe(ctx, addr, h)
}
