package k8s

// clusterType
const (
	TypeKubeIn  = iota // kubernetes cluster in
	TypeKubeOut        // kubernetes cluster out
)

// //TryServiceURL get service url
// func TryServiceURL(ctx context.Context, clusterType int, serviceName string, portName string) string {

// 	var clientset *kubernetes.Clientset
// 	switch clusterType {

// 	case TypeKubeIn:
// 		clientset = TryInClientset()
// 	case TypeKubeOut:
// 		fallthrough
// 	default:
// 		panic(errors.New("don't support type other than TypeKubIn"))
// 	}

// 	service, err := clientset.CoreV1().Services("think").Get(ctx, serviceName, meta_v1.GetOptions{})
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, p := range service.Spec.Ports {

// 		switch p.Name {
// 		case "":
// 			fallthrough //NOTE: 如果没有定义name，默认就是grpc服务端口
// 		case portName: //NOTE:
// 			return serviceName + ":" + strconv.FormatInt(int64(p.Port), 10)
// 		}
// 	}
// 	panic(errors.New("cann't find grpc service port for " + serviceName))
// }

//*****************************in-cluster-client-configuration*******************************//

// var _inclientset *kubernetes.Clientset
// var _inclientsetOnce sync.Once

// // TryInClientset used to create clientset in cluster
// func TryInClientset() *kubernetes.Clientset {

// 	_inclientsetOnce.Do(func() {

// 		config, err := rest.InClusterConfig()
// 		if err != nil {
// 			panic(err)
// 		}
// 		if _inclientset, err = kubernetes.NewForConfig(config); err != nil {
// 			panic(err)
// 		}
// 	})
// 	return _inclientset
// }

// //*****************************out-cluster-client-configuration*******************************//
// var outclientset *kubernetes.Clientset
// // GetOutClientset initialize
// func GetOutClientset() (clientset *kubernetes.Clientset) {

// }
