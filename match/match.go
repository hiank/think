package match

import (
	tg "github.com/hiank/think/net/protobuf/grpc"
)

// Matcher 维护匹配状态
type Matcher struct {

	tokenTeam 		map[string]*Team 		//NOTE: token 与 team 对应map

	// matchedTeams 	map[int]*Team 			//NOTE: 关键字维护的 Team链表，
	waitingTeams 	map[int]*Team 			//NOTE: 等待完成完全匹配的 Team，只会有一个
	waitingFileds 	map[int]*Field 			//NOTE: 等待完成完全匹配的 Field，只会有一个
	// matchedFileds	map

	roleLen 		int 					//NOTE: Team 需要的 Role 的数量
	teamLen 		int 					//NOTE: Field 需要的 Team 的数量

	faCreater		func() FieldAlghorm 	//NOTE: FieldAlghorm 生成器，用于生成一个FieldAlghorm 对象
	responseChan 	chan *tg.Response		//NOTE: 
}

// NewMatcher 创建并初始化一个 Matcher
func NewMatcher(roleLen int, teamLen int, faCreater func() FieldAlghorm, ch chan *tg.Response) *Matcher {

	m := new(Matcher)
	m = &Matcher {

		roleLen: roleLen,
		teamLen: teamLen,
		tokenTeam: make(map[string]*Team),
		waitingTeams: make(map[int]*Team),
		waitingFileds: make(map[int]*Field),
		faCreater: faCreater,
		responseChan: ch,
	}
	return m
}

// Match 将需要匹配的role 加到管理器中，处理
func (m *Matcher) Match(role Role) bool {

	key := role.GetKey()
	team, ok := m.tokenTeam[role.GetId()]
	if !ok {

		if team, ok = m.waitingTeams[key]; !ok {

			team = NewTeam(m.roleLen)
			m.waitingTeams[key] = team
		}
	}
	team.Match(role)		//NOTE: 如果能找到已匹配的team，做匹配处理[已匹配过的role，再调用此api，可能是断线重连了？]
	if !team.Finished() {

		return false
	}
	//NOTE: 以下，为Team 已完全匹配
	field, ok := m.waitingFileds[key]
	if !ok {

		field = NewField(m.teamLen, m.faCreater(), m.responseChan)
		m.waitingFileds[key] = field
	}
	field.Match(team)		//NOTE: 将team匹配到field中
	delete(m.waitingTeams, key)

	ok  = field.Finished()
	if ok {
		delete(m.waitingFileds, key)
	}
	return ok
}

// GetFinishedField 获取token绑定的，已经开始的战局
func (m *Matcher) GetFinishedField(token string) (field *Field, ok bool) {

	team, ok := m.tokenTeam[token]
	if !ok {
		return
	}

	if field, ok = team.GetBindField(); ok {
	
		ok = field.Finished()
		if !ok {
			field = nil
		}
	}
	return
}



// func (m *Matcher) Match() {


// }

// // Matcher 用于维护玩家的配备状态
// type Matcher struct {

// 	nodesLen	int 					//NOTE: Group需要的Node组 的数量
// 	nodeNum 	int 					//NOTE: Group需要的每组Node 的数量

// 	tokenMap 	map[string]*Group		//NOTE: 保存token所在的group
 
// 	matchedMap 	map[int]*Group			//NOTE: 保存已经完全匹配的group
// 	waitingMap 	map[int]*Group			//NOTE: 保存没有完全匹配的group
// }

// // Match 用于对传入的token 对应的玩家进行匹配
// func (m *Matcher) Match(n Node, responseChan chan interface{}) chan interface{} {

// 	if g, ok := m.tokenMap[n.GetToken()]; ok {

// 		return g.GetDataChan()
// 	}

// 	g, ok := m.waitingMap[n.GetKey()]
// 	if !ok {

// 		g = NewGroup(m.nodesLen, m.nodeNum)
// 		m.waitingMap[n.GetKey()] = g
// 	}

// 	err := g.Push(n)
// 	if err != nil {

// 		glog.Errorln("" + err.Error())
// 		return nil
// 	}

// 	if g.Matched() {

		

// 		g.GetMatchedChan() <- true
// 		return g.GetDataChan()
// 	}

// 	<- g.GetMatchedChan()
// 	return g.GetDataChan()
// }

// // Delete 用于删除某个经过匹配处理的玩家
// func (m *Matcher) Delete(token string) {

// 	g, ok := m.tokenMap[token]
// 	if !ok {
// 		return
// 	}
// 	delete(m.tokenMap, token)

// 	// var node *Node
// 	// for _, member := range g.member {

// 	// 	if node, ok := member[token]; ok {
// 	// 		break
// 	// 	}
// 	// }
// 	// if node == nil {
// 	// 	glog.Errorln("no node : " + token)
// 	// 	return
// 	// }
// 	g.Delete(token)


// 	key := node.GetKey()
// 	var gMap map[int]*Group
// 	if g.matched {

// 		gMap = m.matchedMap
// 	} else {

// 		gMap = m.waitingMap
// 	}
// 	gMap[key]
// 	// if g.pre != nil {

// 	// 	g.pre.next = g.next
// 	// } else {			//NOTE: 此处意味着，此节点是初始节点；需要更新匹配链表

// 	// 	node, ok := g.member[token]
// 	// }
// 	// if g.next != nil {

// 	// 	g.next.pre = g.pre
// 	// }

// }


// func (m *Matcher) deleteGroup(g *Group) {

// 	var node *Node
// 	for _, member := range g.member {

// 		if node, ok := member[token]; ok {
// 			break
// 		}
// 	}
// 	if node == nil {
// 		return
// 	}

// 	key := node.GetKey()
// 	// m.matchedMap
// 	var gMap map[int]*Group
// 	if g.matched {

// 		gMap = m.matchedMap
// 	} else {

// 		gMap = m.waitingMap
// 	}

// 	gMap[key]
// }
