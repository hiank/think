package ws_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/hiank/think/net"
	"github.com/hiank/think/net/adapter/ws"
	"github.com/hiank/think/net/box"
	"github.com/hiank/think/net/testdata"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

type tmpAuther string

func (ta tmpAuther) Auth(token string) (uid uint64, err error) {
	arr := strings.Split(token, "_")
	if len(arr) != 2 {
		return 0, fmt.Errorf("invalid token: must format as 'key_number'")
	}
	if arr[0] != string(ta) {
		return 0, fmt.Errorf("invalid token: equal failed")
	}
	return strconv.ParseUint(arr[1], 10, 64)
}

func easyDial(token string) (wc *websocket.Conn, err error) {
	// websocket.NewClient()
	url := &url.URL{Scheme: "ws", Host: "localhost:10240", Path: "/ws"}
	wc, _, err = websocket.DefaultDialer.Dial(url.String(), http.Header{"token": []string{token}})
	return
}

func TestListener(t *testing.T) {
	exit := make(chan bool)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	lis := ws.NewListener(ctx, ws.ListenOption{Addr: ":10240", Auther: tmpAuther("test")})
	go func(t *testing.T) {
		<-exit
		lis.Close()
		_, err := lis.Accept()
		assert.Equal(t, err, io.EOF, err)
		close(exit)
	}(t)

	url := &url.URL{Scheme: "ws", Host: "localhost:10240", Path: "/ws"}
	_, _, err := websocket.DefaultDialer.Dial(url.String(), http.Header{})
	assert.Assert(t, err != nil, "non token info in header")
	//
	// websocket.DefaultDialer.Dial("localhost:30211")
	_, err = easyDial("27")
	assert.Assert(t, err != nil, "invalid token")

	wc, err := easyDial("test_27")
	assert.Equal(t, err, nil, err)

	ic, err := lis.Accept()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, ic.ID, "27", ic.ID)

	// done := make(chan bool)
	go func(ic net.IdentityConn, t *testing.T) {
		m, _ := box.New(&testdata.S_Example{Value: "s-e"})
		err := ic.Send(m)
		assert.Equal(t, err, nil, err)
		m, err = ic.Recv()
		assert.Equal(t, err, nil, err)
		gm, _ := m.GetAny().UnmarshalNew()
		assert.Equal(t, gm.(*testdata.G_Example).GetValue(), "g-v")

		ic.Close()
		// close()
		// close(done)
	}(ic, t)

	mt, buff, err := wc.ReadMessage()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, mt, websocket.BinaryMessage)
	// m := new(box.Message)
	m, err := box.UnmarshalNew[*anypb.Any](buff)
	assert.Equal(t, err, nil, err)
	sm, _ := m.GetAny().UnmarshalNew()
	assert.Equal(t, sm.(*testdata.S_Example).GetValue(), "s-e")

	m, _ = box.New(&testdata.G_Example{Value: "g-v"})
	err = wc.WriteMessage(websocket.BinaryMessage, m.GetBytes())
	assert.Equal(t, err, nil, err)

	// <-done
	mt, _, err = wc.ReadMessage()
	assert.Assert(t, err != nil)
	assert.Equal(t, mt, -1)

	// lis.Close()
	// close(exit)
	exit <- true
	<-exit
}

func TestConn(t *testing.T) {
	exit := make(chan bool)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	lis := ws.NewListener(ctx, ws.ListenOption{Auther: tmpAuther("test1"), Addr: ":10240"})
	go func() {
		<-exit
		lis.Close()
		close(exit)
	}()

	wc, err := easyDial("test1_11")
	assert.Equal(t, err, nil, err)
	c := ws.Export_newConn(wc)
	m, _ := box.New(&testdata.P_Example{Value: "p-v"})
	err = c.Send(m)
	assert.Equal(t, err, nil, err)

	sc, err := lis.Accept()
	assert.Equal(t, err, nil, err)
	sm, _ := sc.Recv()
	sv, _ := sm.GetAny().UnmarshalNew()
	assert.Equal(t, sv.(*testdata.P_Example).GetValue(), "p-v")

	m, _ = box.New(&testdata.MessageTest1{Key: "m-t"})
	err = sc.Send(m)
	assert.Equal(t, err, nil, err)

	m, err = c.Recv()
	assert.Equal(t, err, nil, err)
	v1, _ := m.GetAny().UnmarshalNew()
	assert.Equal(t, v1.(*testdata.MessageTest1).GetKey(), "m-t")

	err = c.Close()
	assert.Equal(t, err, nil, err)

	_, err = sc.Recv()
	assert.Assert(t, err != nil)

	exit <- true
	<-exit
}

func TestWithDefaultListenOption(t *testing.T) {
	opt := ws.Export_withDefaultListenOption(ws.ListenOption{Addr: "11"})
	uid, err := opt.Auther.Auth("")
	assert.Equal(t, uid, uint64(0))
	assert.Equal(t, err, ws.ErrUnimplementedAuther)
	assert.Equal(t, opt.Addr, "11")

	opt = ws.Export_withDefaultListenOption(ws.ListenOption{Auther: tmpAuther("tt")})
	uid, err = opt.Auther.Auth("tt_101")
	assert.Equal(t, uid, uint64(101))
	assert.Equal(t, err, nil, err)
	assert.Equal(t, opt.Addr, "")
}
