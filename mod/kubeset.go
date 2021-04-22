package mod

import (
	"context"

	"github.com/hiank/think"
	"github.com/hiank/think/net/k8s"
	"k8s.io/client-go/kubernetes"
)

var (
	KubesetIn = &kubeset{}
)

type kubeset struct {
	clientset *kubernetes.Clientset
	think.IgnoreDepend
	think.IgnoreOnStart
	think.IgnoreOnStop
	think.IgnoreOnDestroy
}

//OnCreate 此阶段，需要把配置数据注册到ConfigMod
func (set *kubeset) OnCreate(ctx context.Context) (err error) {
	set.clientset, err = k8s.NewInClientset()
	return
}

//AutoAddr 自动生成addr，会从clientset中根据端口名获取端口号
func (set *kubeset) AutoAddr(svcName, portName string) (string, error) {
	return k8s.AutoAddr(context.Background(), set.clientset, svcName, portName)
}
