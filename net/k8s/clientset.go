package k8s

import (
	"context"
	"errors"
	"strconv"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//NewInClientset 创建一个新的集群内Clientset
func NewInClientset() (clientset *kubernetes.Clientset, err error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		clientset, err = kubernetes.NewForConfig(config)
	}
	return
}

//AutoAddr 自动生成addr，根据端口名从k8s中读到端口号
func AutoAddr(ctx context.Context, clientset *kubernetes.Clientset, svcName, portName string) (string, error) {
	svc, err := clientset.CoreV1().Services("think").Get(ctx, svcName, v1.GetOptions{})
	if err != nil {
		return "", err
	}
	for _, p := range svc.Spec.Ports {
		switch p.Name {
		case "":
			fallthrough //NOTE: 如果没有定义name，默认就是grpc服务端口
		case portName: //NOTE:
			return svcName + ":" + strconv.FormatInt(int64(p.Port), 10), nil
		}
	}
	return "", errors.New("cannot find port named " + portName)
}
