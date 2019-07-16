package k8s

import (
	"bytes"
	"time"
	"context"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc"
	"strconv"

	"errors"
	"github.com/golang/glog"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hiank/think/util"
)
	

// clusterType 
const (

	TypeKubIn 	= iota 	// kubernetes cluster in
	TypeKubOut			// kubernetes cluster out
) 

//ServiceNameWithPort 通过消息名 得到服务名与端口号连接字符串
func ServiceNameWithPort(clusterType int, serviceName string, portName string) (addr string, err error) {

	var clientset *kubernetes.Clientset
	switch (clusterType) {

	case TypeKubIn:
		clientset = GetInClientset()
	case TypeKubOut: fallthrough
	default:
		err = errors.New("don't support type other than TypeKubIn")
		return
	}

	// serviceName := msgName + "-service"
	service, err := clientset.CoreV1().Services("think").Get(serviceName, meta_v1.GetOptions{})
	if err != nil {
		glog.Error("cann't get service named : " + serviceName + " : " + err.Error())
		return
	}

	var port int32
	L: for _, p := range service.Spec.Ports {

		switch p.Name {
		case "": fallthrough		//NOTE: 如果没有定义name，默认就是grpc服务端口
		case portName:				//NOTE:
			port = p.Port
			break L
		}
	}

	if port == 0 {					//NOTE: 如果没有找到端口
		err = errors.New("cann't find grpc service port for " + serviceName)
		glog.Error(err)
		return
	}

	buffer := bytes.NewBufferString(serviceName)
	buffer.WriteByte(':')
	buffer.WriteString(strconv.FormatInt(int64(port), 10))

	addr = buffer.String()
	return
}

// DailToCluster 连接到k8s
func DailToCluster(clusterType int, msgName string) (*grpc.ClientConn, error) {

	addr, err := ServiceNameWithPort(clusterType, msgName + "-service", "grpc")
	if err != nil {
		return nil, err
	}
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		glog.Error(err)
		return cc, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	for {
		s := cc.GetState()
		if s == connectivity.Ready {
			break
		}
		if !cc.WaitForStateChange(ctx, s) {
			err = ctx.Err()
			cc = nil
			glog.Error(err, ":", s.String())
			break
		}
	}
	return cc, err
}


//*****************************in-cluster-client-configuration*******************************//

var inclientset *kubernetes.Clientset
// GetInClientset used to create clientset in cluster
func GetInClientset() (clientset *kubernetes.Clientset) {

	// glog.Infof("do get in clientset : %v\n", inclientset)
	if inclientset == nil {

		defer util.RecoverErr("GetInClientset error : ")

		config, err := rest.InClusterConfig()
		util.PanicErr(err)

		// creates the clientset
		inclientset, err = kubernetes.NewForConfig(config)
		util.PanicErr(err)

	}
	clientset = inclientset
	return
}


// //*****************************out-cluster-client-configuration*******************************//
// var outclientset *kubernetes.Clientset
// // GetOutClientset initialize 
// func GetOutClientset() (clientset *kubernetes.Clientset) {


// }
