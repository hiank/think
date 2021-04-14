package k8s

// import (
// 	"context"
// 	"errors"

// 	"github.com/hiank/think"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/rest"
// )

// type modClientset struct {
// 	kubeType  int
// 	clientset *kubernetes.Clientset
// 	think.IgnoreDepend
// 	think.IgnoreOnDestroy
// 	think.IgnoreOnStart
// 	think.IgnoreOnStop
// }

// func (mcs *modClientset) OnCreate(ctx context.Context) (err error) {
// 	var config *rest.Config
// 	switch mcs.kubeType {
// 	case TypeKubeIn:
// 		if config, err = rest.InClusterConfig(); err != nil {
// 			return
// 		}
// 	case TypeKubeOut:
// 		return errors.New("not support out cluster current")
// 	}

// 	mcs.clientset, err = kubernetes.NewForConfig(config)
// 	return
// }

// // func (mcs *modClientset)

// var InClientset = &modClientset{kubeType: TypeKubeIn}

// var OutClientset = &modClientset{kubeType: TypeKubeOut}
