//internal tests
package pool

import (
	"strconv"
	"testing"

	"github.com/hiank/think/pb"
	"github.com/hiank/think/token"
)

type testIO func(*pb.Message) error

func (ti testIO) Recv() (*pb.Message, error) {

	return nil, nil
}

func (ti testIO) Send(msg *pb.Message) error {

	return ti(msg)
}

var tokenInt = 0

func newTestConn(arr ...IO) *Conn {

	var handler IO
	if len(arr) > 0 {
		handler = arr[0]
	} else {
		handler = testIO(func(*pb.Message) error { return nil })
	}

	tokenInt++
	tok := token.GetBuilder().Get(strconv.Itoa(tokenInt))
	return newConn(tok, handler)
}

//***** test connhub ****//

// func TestConnHubAdd(t *testing.T) {

// 	hub, c := newConnHub(context.Background()), newTestConn()
// 	r := &req{tag: typeAdd, param: c, res: make(chan interface{})}
// 	hub.req <- r
// 	success, ok := <-r.res
// 	assert.Equal(t, ok, true)
// 	assert.Equal(t, success.(bool), true)
// 	assert.Equal(t, len(hub.hub), 1)
// }

// func TestConnHubDel(t *testing.T) {

// 	hub := newConnHub(context.Background())
// 	del := &req{tag: typeDel, param: strconv.Itoa(tokenInt), res: make(chan interface{})}

// 	hub.req <- del
// 	success, ok := <-del.res
// 	assert.Equal(t, ok, true)
// 	assert.Equal(t, success.(bool), true)
// }

// func TestConnHubFind(t *testing.T) {

// 	hub := newConnHub(context.Background())
// 	findReq := &req{tag: typeFind, param: strconv.Itoa(tokenInt), res: make(chan interface{})}

// 	var hopeC *conn
// 	hub.req <- findReq
// 	c, ok := <-findReq.res
// 	assert.Equal(t, ok, true)
// 	assert.Equal(t, c.(*conn), hopeC)

// 	hopeC = newTestConn()
// 	addReq := &req{tag: typeAdd, param: hopeC, res: make(chan interface{})}
// 	hub.req <- addReq
// 	<-addReq.res

// 	findReq.param = hopeC.ToString()
// 	hub.req <- findReq
// 	c, ok = <-findReq.res
// 	assert.Equal(t, ok, true)
// 	assert.Equal(t, c.(*conn), hopeC)
// }

// func TestConnHubSend(t *testing.T) {

// 	c := newTestConn(testIO(func(pbMsg *pb.Message) error {
// 		return nil
// 	}))
// 	hub, msg := newConnHub(context.Background()), NewMessage(&pb.Message{Token: c.ToString()}, c.Token.Derive())
// 	sendReq := &req{tag: typeSend, param: msg, res: make(chan interface{})}

// 	hub.req <- sendReq
// 	err, ok := <-sendReq.res
// 	// var hopeErr error
// 	assert.Equal(t, ok, true)
// 	assert.Error(t, err.(error), "cann't find conn tokened "+msg.ToString())

// 	addReq := &req{tag: typeAdd, param: c, res: make(chan interface{})}
// 	hub.req <- addReq
// 	<-addReq.res

// 	hub.req <- sendReq
// 	err, ok = <-sendReq.res
// 	assert.Equal(t, ok, true)
// 	assert.Equal(t, err, nil)
// }

//***** test conn *****//

func TestConnListen(t *testing.T) {

}

func TestConnWork(t *testing.T) {

}
