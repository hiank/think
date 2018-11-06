package match

// import (
// 	"github.com/golang/glog"
// 	"errors"
// )

// // Group 用于维护一个战局中的玩家
// type Group struct {

// 	matchedChan chan bool 					//NOTE: 匹配完全后，将写true 入此chan
// 	dataChan 	chan interface{}			//NOTE: 用于战局接受request 信息

// 	matched 	bool						//NOTE: 用于记录当前Group 是否匹配完整
// 	members 	[]map[string]Node			//NOTE: 玩家组 map[token]chan
// 	memberNum 	int							//NOTE: 每组玩家的数量

// 	field 		Field 						//NOTE: 绑定战局对象；具体战场算法

// 	pre 		*Group 						//NOTE: 上一组玩家；方便链表增删
// 	next 		*Group						//NOTE: 下一组玩家；因为战局随时可能终结，用链表维护，增减更高效，方便
// }

// // NewGroup 创建一个新的Group
// func NewGroup(len int, num int) *Group {

// 	g := new(Group)
// 	g = &Group {

// 		matchedChan: make(chan bool),
// 		dataChan: make(chan interface{}),
// 		matched: false,
// 		members: make([]map[string]Node, len),
// 		memberNum: num,
// 		pre: nil,
// 		next: nil,
// 	}
// 	return g
// }

// // GetDataChan 用于战局接受request 数据
// func (g *Group) GetDataChan() chan interface{} {

// 	return g.dataChan
// }


// // GetMatchedChan 匹配成功后会向此chan谢 true
// func (g *Group) GetMatchedChan() chan bool {

// 	return g.matchedChan
// }

// // Matched 判断此Group是否已完全匹配
// func (g *Group) Matched() bool {

// 	return g.matched
// }


// // BindField 用于将战局相关的算法绑定到Group中
// func (g *Group) BindField(field Field) {

// 	g.field = field
// }


// // Push 向玩家组中无差别添加 单元
// func (g *Group) Push(n Node) (err error) {

// 	num := 0
// 	success := false
// 	for _, member := range g.members {

// 		memberLen := len(member)
// 		num += memberLen
// 		if memberLen < g.memberNum {

// 			member[n.GetToken()] = n
// 			success = true
// 			num++
// 			break
// 		}
// 	}

// 	g.matched = (num == (len(g.members) * g.memberNum))
// 	if !success {
// 		err = errors.New("Group full")
// 	}
// 	return
// }

// // Delete 删除token对应的Node
// func (g *Group) Delete(token string) {

// 	success, empty := false, true
// 	for _, member := range g.members {

// 		if _, ok := member[token]; ok {

// 			delete(member, token)
// 			success = true
// 			break
// 		}
// 	}

// 	if !success {
// 		glog.Warningln("Group Delete no Node with token : " + token)
// 	}


// }
