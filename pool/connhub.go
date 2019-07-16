package pool

import (
	"github.com/golang/glog"
	"github.com/hiank/think/pb"
	"container/list"
	"sync"
	"errors"
)

//ConnHub used to maintain conn
type ConnHub struct {

	mtx sync.RWMutex				//NOTE: 读写锁，ConnHub 会在不同goroutine中添删conn
	hub map[string]*list.List		//NOTE: 用于保存conn，key 一般为服务名
}

//newConnHub 创建一个新的ConnHub
func newConnHub() *ConnHub {

	return &ConnHub {

		hub : make(map[string]*list.List),
	}
}


//CheckConnected 检查Conn 是否已连接
func (ch *ConnHub) CheckConnected(key, token string) bool {

	// key, token := it.GetKey(), it.GetToken()
	connected := false
	if queue, ok := ch.hub[key]; ok {

		for element := queue.Front(); element != nil; element = element.Next() {

			if token == element.Value.(*Conn).GetToken().ToString() {

				connected = true
				break
			}
		}
	}
	return connected
}

//Push 将新的Conn 加入到队列尾
func (ch *ConnHub) Push(conn *Conn) {

	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	key := conn.GetKey()
	glog.Infoln("ConnHub Push key : ", key)
	queue, ok := ch.hub[key]
	if !ok {
		queue = list.New()
		ch.hub[key] = queue
	}
	conn.Element = queue.PushBack(conn)
	conn.Update()
}

//Handle 处理数据发送，总感觉这边会有性能问题，如果有超级多的玩家同时在线，比如1000万，每次要发送一个消息都要遍历查找一遍，可能会卡死
func (ch *ConnHub) Handle(msg *pb.Message) (err error) {

	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	key := msg.GetKey()
	queue, ok := ch.hub[key]
	glog.Infoln("connhub handle keyed : ", key)
	if !ok {
		err = errors.New("connhub has no list keyed " + key)
		glog.Infoln(err)
		return
	}
	token := msg.GetToken()
	var conn *Conn
	for element := queue.Front(); element != nil; element = element.Next() {

		c := element.Value.(*Conn)
		if token == c.GetToken().ToString() {
			conn = c
			break
		}
	}
	if conn == nil {

		err = errors.New("connhub has no conn tokened " + token)
		glog.Infoln(err)
		return
	}
	err = conn.Send(msg)
	return
}

//Update 每次有通讯，将conn 移到队尾，提高清理效率
func (ch *ConnHub) Update(conn *Conn) {

	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	key := conn.GetKey()
	glog.Infof("updage conn keyed %s, tokened %s\n", key, conn.GetToken().ToString())
	ch.hub[key].MoveToBack(conn.Element)
}


//Upgrade 清理超时连接
func (ch *ConnHub) Upgrade() {

	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	for _, v := range ch.hub {

		ch.upgrade(v)
	}
}

func (ch *ConnHub) upgrade(queue *list.List) {

	element := queue.Front()
	for element != nil {

		if !element.Value.(Timer).TimeOut() {
			break
		}

		glog.Infoln("conn keyed ", element.Value.(*Conn).GetKey(), " timeout !")

		cur := element
		element = element.Next()
		cur.Value.(*Conn).GetToken().Cancel()			//NOTE: 关闭Context，将触发pool Listen中调用ConnHub Remove的逻辑，conn将在Remove中被清除
	}
}

//Remove 当conn 关闭的时候，调用这个方法，清除 
//NOTE: 这个方法只在pool中调用
func (ch *ConnHub) Remove(conn *Conn) {

	ch.mtx.Lock()

	ch.hub[conn.GetKey()].Remove(conn.Element)
	conn.Element = nil

	ch.mtx.Unlock()
}