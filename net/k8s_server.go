package net

import (
	"errors"
	"github.com/hiank/think/net/k8s"
)

//ServeK8s 启动一个k8s服务
func ServeK8s(addr string, h k8s.MessageHandler) error {

	if netCtx == nil {
		return errors.New("net.Init should be called first")
	}
	return k8s.ListenAndServe(netCtx, addr, h)
}
