package net

import (
	"container/list"
	"time"
	"io"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gorilla/websocket"
	gp "github.com/hiank/think/net/protobuf/grpc"

	proto "github.com/hiank/think/net/protobuf"
)


// Conn is the base data type operate of net
type Conn struct {

	unix 		int64 				//NOTE: 最近收到消息的时间戳，用于判断超时
	conn 		*websocket.Conn
	element 	*list.Element 		//NOTE: 保存在ConnPool/connlist 中的Element，方便更新时间戳后更新位置

	datachan 	chan *any.Any
	token 		string
}

// NewConn create an new Conn for server
func NewConn(conn *websocket.Conn, token string) *Conn {

	c := &Conn{

		unix 		: time.Now().Unix(),
		conn		: conn,
		datachan	: make(chan *any.Any),
		token 		: token,
	}
	return c
}


// ReadMessage read from conn and analyze it
func (c *Conn) ReadMessage() (data *gp.Request, err error) {

	c.unix = time.Now().Unix()				//NOTE: 更新时间戳

	_, buf, err := c.conn.ReadMessage()
	if err != nil {
		return
	}

	msg, err := proto.PBDecode(buf)
	if err != nil {
		return
	}

	data = &gp.Request{Token:[]byte(c.token), Data:msg}
	return
}

// SendAsync 异步发送，收到k8s 返回的数据，处理后转发到客户端
func (c *Conn) SendAsync(quit chan bool) {

L:	for {
		select {
		case <-quit: break L
		case msg := <-c.datachan:
			d, err := proto.PBEncode(msg)
			if err != nil {
				continue L
			}
			err = c.conn.WriteMessage(websocket.BinaryMessage, d)
			if err == io.EOF {
				quit <- true
				break L
			}	
		}
	}
}

// SetElement 保存ConnPool/listconn 中的Element
func (c *Conn) SetElement(e *list.Element) {

	c.element = e
}

// GetElement 获得element，用于更新element的位置
func (c *Conn) GetElement() *list.Element {

	return c.element
}

// GetDatachan to get schan
func (c *Conn) GetDatachan() chan *any.Any {

	return c.datachan
}

// GetToken used to get token of Conn
func (c *Conn) GetToken() string {

	return c.token
}

// Close 关闭此连接
func (c *Conn) Close() {

	c.conn.Close()
	close(c.datachan)
	c.element = nil
}

// TimeOut 判断连接是否超时，第一次超时尝试建立通信
func (c *Conn) TimeOut() (out bool) {

	t := time.Now()
	interval := t.Unix() - c.unix

	out = (interval > 10)		//NOTE: 设定为10秒，实际需要配表
	return
}


// ConnPool used to manager wserver's conn
// 处理conn的增删改查
type ConnPool struct {

	connaddchan 	chan *Conn				//NOTE: add conn
	conndelchan 	chan *Conn				//NOTE: delete conn
	connupdatechan 	chan *Conn				//NOTE: conn收到消息后，更新此chan，用于改变element的位置

	tokenconn 		map[string]*Conn		//NOTE: key is token
	connlist 		*list.List				//NOTE: 保存conn 的链表
	waitlist 		*list.List 				//NOTE: 保存未传入token 的conn 链表
}

var connpool *ConnPool
// GetConnPool get initialized static value
func GetConnPool() (cm *ConnPool) {

	if connpool != nil {
		cm = connpool
		return
	}

	cm = &ConnPool {

		connaddchan: make(chan *Conn),
		conndelchan: make(chan *Conn),
		connupdatechan: make(chan *Conn),
		tokenconn: make(map[string]*Conn),
		connlist: list.New(),
	}
	go cm.loop()
	connpool = cm
	return
}

// ReleaseConnPool clean the static object
func ReleaseConnPool() {

	if connpool == nil {
		return
	}

	head := connpool.connlist.Front()
	for head != nil {

		head.Value.(*Conn).Close()
		next := head.Next()
		connpool.connlist.Remove(head)
		head = next
	}
	connpool.tokenconn = nil
	connpool.connlist = nil

	close(connpool.connaddchan)
	close(connpool.conndelchan)
	close(connpool.connupdatechan)
	connpool = nil
}


// // GetConnaddChan used to get conn chan
// func (cp *ConnPool) GetConnaddChan() chan *Conn {

// 	return cp.connaddchan
// }

// // GetConndelChan get conndelchan to post delete message for conn
// func (cp *ConnPool) GetConndelChan() chan *Conn {

// 	return cp.conndelchan
// }

// // GetConnupdateChan get connupdatechan
// func (cp *ConnPool) GetConnupdateChan() chan *Conn {

// 	return cp.connupdatechan
// }

// AddConn 向ConnPool中添加conn
func (cp *ConnPool) AddConn(conn *Conn) {

	cp.connaddchan <- conn
}

// DelConn 从ConnPool中删除conn
func (cp *ConnPool) DelConn(conn *Conn) {

	cp.conndelchan <- conn
}

// UpdateConn 更新conn的连接状态
func (cp *ConnPool) UpdateConn(conn *Conn) {

	cp.connupdatechan <- conn
}

func (cp *ConnPool) loop() {

	for {

		select {
		case conn := <-cp.connaddchan:		//NOTE: 收到新建立的连接，加入管理列表
			cp.add(conn)
		case conn := <-cp.conndelchan:		//NOTE: 删除连接
			cp.del(conn)
		case conn := <-cp.connupdatechan:	//NOTE: 更新conn在connlist中的位置
			cp.connlist.MoveToBack(conn.GetElement())
		case <-time.After(time.Second * 10):	//NOTE: 每隔10秒[配表]，执行一次清理
			cp.clean()
		}
	}
}


// addConn 加入c，此时，还没有token
func (cp *ConnPool) add(c *Conn) {

	token := string(c.GetToken())		//NOTE: 断线重连，token会一样的吗？

	if oldConn, ok := cp.tokenconn[token]; ok {

		cp.del(oldConn)
	}

	element := cp.connlist.PushBack(c)
	c.SetElement(element)
	
	cp.tokenconn[token] = c
}

// delConn 删除指定conn
func (cp *ConnPool) del(c *Conn) {


	key := string(c.token)
	if _, ok := connpool.tokenconn[key]; ok {

		delete(connpool.tokenconn, key)
	}

	element := c.GetElement()
	c.Close()
	cp.connlist.Remove(element)
}

// clean 执行清理操作，对超时的操作conn 进行处理[判断连接是否还有效]
func (cp *ConnPool) clean() {

	head := cp.connlist.Front()
	for head != nil && head.Value.(*Conn).TimeOut() {

		//NOTE: 此处执行超时操作

		head = head.Next()
	}
}


// RecvAsync run in deferent goroutine for recv message
func (cp *ConnPool) RecvAsync(quit chan bool, k8schan chan K8sResponse) {

L:	for {
		select {
		case <-quit: return
		case response := <-k8schan:		//NOTE: 从k8s 服务器收到消息并处理

			data := response.GetData()
			tokens := response.GetTokens()
			if len(tokens) == 0 {		//NOTE: 如果tokens 没有值，表明需要全体在线广播

				for _, conn := range cp.tokenconn {

					conn.GetDatachan() <- data
				}
				continue L
			}
			for token := range tokens {

				if conn, ok := cp.tokenconn[string(token)]; ok {

					conn.GetDatachan() <- data
				}
			}
		}
	}
}


