package match

import (
	"github.com/golang/glog"
	// tg "github.com/hiank/think/net/protobuf/grpc"
)

// Team 一组合作的玩家，玩家数量最少为1
type Team struct {

	list 			[]Role							//NOTE: 保存匹配的玩家，考虑到玩家数量比较有限，使用数组存储
	cap 			int 							//NOTE: list 的cap，为Role的需要数量

	// responseMap		map[chan *tg.Response][]string	//NOTE: 保存response chan 对应的id 

	bindField 		*Field 							//NOTE: 所在的field
	// pre 		*Team		//NOTE: 链表上一元素
	// next 		*Team 		//NOTE: 链表下一元素
}


// NewTeam 创建一个新的Team，并初始化 
func NewTeam(roleLen int) *Team {

	team := new(Team)
	team = &Team {
		list: make([]Role, 0, roleLen),
		// matched: false,
		cap: roleLen,
	}
	return team
}


// Match 将role 匹配到Team中
func (t *Team) Match(role Role) {

	id := role.GetId()
	for idx, r := range t.list {

		if r.GetId() == id {

			// t.deleteResponse(r)
			t.list[idx] = role
			// t.appendResponse(role)
			return
		}
	}
	if t.Finished() {

		glog.Warningln("team is matched")
		return
	}

	t.list = append(t.list, role)
	// t.appendResponse(role)
}

// GetRoles 获得当前team中的role列表
func (t *Team) GetRoles() []Role {

	return t.list
}

// func (t *Team) appendResponse(role Role) {

// 	ch := role.GetResponseChan()
// 	list, ok := t.responseMap[ch]
// 	if !ok {

// 		list = make([]string, 0, t.cap)
// 		// t.responseMap[ch] = list
// 	}

// 	roleID := role.GetId()
// 	for _, id := range list {

// 		if roleID == id {

// 			return
// 		}
// 	}
// 	t.responseMap[ch] = append(list, roleID)
// }

// func (t *Team) deleteResponse(role Role) {

// 	ch := role.GetResponseChan()
// 	list, ok := t.responseMap[ch]
// 	if !ok {
// 		return
// 	}

// 	roleID, idx := role.GetId(), -1
// 	for i, id := range list {

// 		if roleID == id {

// 			idx = i
// 			break
// 		}
// 	}

// 	switch {

// 	case idx == -1:
// 	case len(list) == 1:
// 		delete(t.responseMap, ch)
// 	default:
// 		t.responseMap[ch] = append(list[:idx], list[idx+1:]...)
// 	}

// }

func (t *Team) Finished() bool {

	return t.cap == len(t.list)
}


func (t *Team) BindField(field *Field) {

	t.bindField = field
	// field.AddResponseInfo(t.responseMap)
}

// GetBindField 获得绑定的Field
func (t *Team) GetBindField() (field *Field, ok bool) {


	return
}

