package k8s

import (
	"context"
	"errors"
	"sync"

	"github.com/golang/glog"
	"github.com/hiank/think/utils"
	"github.com/hiank/think/utils/robust"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)
	

// clusterType 
const (

	TypeKubIn 	= iota 	// kubernetes cluster in
	TypeKubOut			// kubernetes cluster out
) 

//ServiceNameWithPort 通过消息名 得到服务名与端口号连接字符串
func ServiceNameWithPort(ctx context.Context, clusterType int, serviceName string, portName string) (addrWithPort string, err error) {

	var clientset *kubernetes.Clientset
	switch (clusterType) {

	case TypeKubIn:
		clientset = GetInClientset()
	case TypeKubOut: fallthrough
	default:
		err = errors.New("don't support type other than TypeKubIn")
		return
	}

	service, err := clientset.CoreV1().Services("think").Get(ctx, serviceName, meta_v1.GetOptions{})
	if err != nil {
		glog.Error("cann't get service named : " + serviceName + " : " + err.Error())
		return
	}

	var port uint16
	L: for _, p := range service.Spec.Ports {

		switch p.Name {
		case "": fallthrough		//NOTE: 如果没有定义name，默认就是grpc服务端口
		case portName:				//NOTE:
			port = uint16(p.Port)
			break L
		}
	}

	if port == 0 {					//NOTE: 如果没有找到端口
		err = errors.New("cann't find grpc service port for " + serviceName)
		glog.Error(err)
		return
	}
	addrWithPort = utils.WithPort(serviceName, port)
	return
}


//*****************************in-cluster-client-configuration*******************************//

var _inclientset *kubernetes.Clientset
var _inclientsetOnce sync.Once

// GetInClientset used to create clientset in cluster
func GetInClientset() *kubernetes.Clientset {

	_inclientsetOnce.Do(func ()  {
		
		defer robust.Recover(robust.Warning)
		
		config, err := rest.InClusterConfig()
		robust.Panic(err)

		// creates the clientset
		_inclientset, err = kubernetes.NewForConfig(config)
		robust.Panic(err)
	})
	return _inclientset
}


// //*****************************out-cluster-client-configuration*******************************//
// var outclientset *kubernetes.Clientset
// // GetOutClientset initialize 
// func GetOutClientset() (clientset *kubernetes.Clientset) {


// }
