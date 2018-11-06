package match

import (
	"github.com/golang/glog"
	tg "github.com/hiank/think/net/protobuf/grpc"
	"github.com/golang/protobuf/ptypes/any"
)

// FieldAlghorm 战斗算法
type FieldAlghorm interface {

	GetEndChan() chan bool 			//NOTE: 用于外部监听战斗结束
	Do(*any.Any) *any.Any 			//NOTE: 执行client发来的操作指令
}


// Field 地图
type Field struct {

	list 			[]*Team 			//NOTE: 保存team
	cap 			int
	datachan		chan *tg.Request	//NOTE: 在战斗中监听client 操作
	responseChan 	chan *tg.Response 	//NOTE: 战斗状态同步
	// responseMap map[chan *tg.Response][]string 	//NOTE: 保存需要广播的chan
	tokens 			[][]byte			//NOTE: 用于保存当前Field中所有的玩家的token，用于广播发送消息

	althorm 		FieldAlghorm		//NOTE: 具体的战斗算法
}

// NewField 创建一个新的Field
func NewField(teamLen int, althorm FieldAlghorm, ch chan *tg.Response) *Field {

	f := new(Field)
	f = &Field {
		list: make([]*Team, 0, teamLen),
		cap: teamLen,
		datachan: make(chan *tg.Request),
		responseChan: ch,
		althorm: althorm,
	}
	return f
}

// Match 匹配逻辑
func (f *Field) Match(t *Team) {

	if f.Finished() {

		glog.Errorln("filed is matched")
		return
	}

	roles := t.GetRoles()	//NOTE: 此处的每个role 的token 一定不同
	if f.tokens == nil {

		f.tokens = make([][]byte, 0, f.cap * len(roles))
	}
	for _, role := range roles {

		f.tokens = append(f.tokens, []byte(role.GetToken()))
	}

	t.BindField(f)
	f.list = append(f.list, t)
	if f.Finished() {

		go f.run()
	}
}

// Finished 判断匹配是否已完全
func (f *Field) Finished() bool {

	return f.cap == len(f.list)
}

// run 启动战斗
func (f *Field) run() {

	for {
		select {
		case <-f.althorm.GetEndChan(): return
		case data := <-f.datachan: 
			msg := f.althorm.Do(data.GetData())
			f.responseChan <- &tg.Response{Tokens:f.tokens, Data:msg}
		}
	}
}

// Request 发送请求到协程中
func (f *Field) Request(tank *tg.Request) {

	f.datachan <- tank
}
