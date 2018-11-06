package cluster

import (
	"github.com/golang/glog"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hiank/think/util"

)


// GetIPAndPortInKub get ip and port for msgName in clientset
func GetIPAndPortInKub(clientset *kubernetes.Clientset, msgName string) (ip string, port int32, err error) {

	pod, err := clientset.CoreV1().Pods("think").Get(msgName, meta_v1.GetOptions{})
	if err != nil {
		// panic(err)
		glog.Error("no pod for msg " + msgName + ":" + err.Error())
		return
	}

	ip = pod.Status.PodIP
	port = pod.Spec.Containers[0].Ports[0].ContainerPort
	return
}


//*****************************in-cluster-client-configuration*******************************//

var inclientset *kubernetes.Clientset
// GetInClientset used to create clientset in cluster
func GetInClientset() (clientset *kubernetes.Clientset) {

	if inclientset == nil {

		defer util.RecoverErr("GetInClientset error : ")

		config, err := rest.InClusterConfig()
		util.PanicErr(err)
	
		// creates the clientset
		clientset, err = kubernetes.NewForConfig(config)
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