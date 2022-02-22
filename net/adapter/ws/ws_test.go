package ws_test

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hiank/think/net/adapter/ws"
	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/net/testdata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

func TestSendChan(t *testing.T) {
	t.Run("closed chan", func(t *testing.T) {
		ch := make(chan bool)
		close(ch)

		defer func(t *testing.T) {
			r := recover()
			assert.Assert(t, r != nil)
			assert.Equal(t, r.(error).Error(), "send on closed channel")
		}(t)
		ch <- true
	})

	t.Run("nil chan", func(t *testing.T) {
		var ch chan bool
		select {
		case ch <- true:
			assert.Assert(t, false)
		default:
			assert.Assert(t, true)
		}
	})
}

type testStorage struct {
}

func (ts *testStorage) Auth(token string) (uid uint64, err error) {
	uid, err = strconv.ParseUint(token, 10, 64)
	return
}

func TestListener(t *testing.T) {
	t.Run("new-close", func(t *testing.T) {
		ts := &testStorage{}
		uid, _ := ts.Auth("11")
		assert.Equal(t, uid, uint64(11))
		l := ws.NewListener(ts, ":10240")
		l.Close()

		<-time.After(time.Millisecond) //NOTE: wait for server listener stopped

		l.Close()
	})

	// <-time.After(time.Millisecond)

	// assert.
}

func TestConn(t *testing.T) {

	l := ws.NewListener(&testStorage{}, ":10240")
	defer l.Close()

	// websocket.NewClient()
	url := &url.URL{Scheme: "ws", Host: "localhost:10240", Path: "/ws"}
	_, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	// t.Log(err)
	assert.Assert(t, err != nil)
	// assert.Assert(t, resp == nil)
	// t.Log(resp.StatusCode, )
	// assert.Equal(t, resp.StatusCode, http.StatusNonAuthoritativeInfo)

	_, resp, _ := websocket.DefaultDialer.Dial(url.String(), http.Header{"token": []string{"not number"}})
	// assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusUnauthorized, resp)

	cliConn, _, err := websocket.DefaultDialer.Dial(url.String(), http.Header{"token": []string{"11"}})
	assert.Assert(t, err == nil)

	srvConn, err := l.Accept()
	assert.Assert(t, err == nil)

	wait := make(chan bool)
	// srvConn.Read()
	go func(t *testing.T) {
		d, _ := srvConn.Recv()
		amsg := new(anypb.Any)
		proto.Unmarshal(d.Bytes(), amsg)
		msg, _ := amsg.UnmarshalNew()
		assert.Equal(t, msg.(*testdata.AnyTest1).GetName(), "ll")
		// close(wait)
		wait <- true
	}(t)

	any, _ := anypb.New(&testdata.AnyTest1{Name: "ll"})
	b, _ := proto.Marshal(any)
	err = cliConn.WriteMessage(websocket.BinaryMessage, b)
	assert.Assert(t, err == nil)

	<-wait

	any, _ = anypb.New(&testdata.AnyTest2{Hope: "hh"})
	b, _ = proto.Marshal(any)
	doc, _ := pb.MakeM(b)
	srvConn.Send(doc)

	_, b, _ = cliConn.ReadMessage()
	// var msg anypb.Any
	proto.Unmarshal(b, any)
	msg, _ := any.UnmarshalNew()
	assert.Equal(t, msg.(*testdata.AnyTest2).GetHope(), "hh")

	// assert.Equal(t, srvConn.GetIdentity(), uint64(11))

	srvConn.Close()
	_, _, err = cliConn.ReadMessage()
	assert.Assert(t, err != io.EOF)
}
