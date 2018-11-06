package match

// // Node 用于保存玩家信息；初步设想到redis中获取; 结构另定
// type Node struct {

// 	ch chan *tg.Response
// 	cups		int 		//NOTE: 奖杯数，用于匹配 
// }

// // GetKey 用于获取匹配关键字
// func (n *Node) GetKey() int {

// 	return n.cups
// }

// // GetResponseChan 用于获得 通知的 chan
// func (n *Node) GetResponseChan() chan *tg.Response {

// 	return n.ch
// }

// Node 战局中最小单位
type Node interface {

	GetToken() string 						//NOTE: 用于获取 最小单位对应的token，服务中唯一
	GetKey() int 							//NOTE: 用于获取匹配关键字
	GetResponseChan() chan interface{}		//NOTE: 用于获取 通知战局 的chan
}

