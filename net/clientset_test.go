package net_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/box"
	"github.com/hiank/think/net/one"
	"github.com/hiank/think/net/testdata"
	"gotest.tools/v3/assert"
)

type tmpKnower struct {
	invalidKey string
}

// func (tmpKnower) Identity(uid string) (id string, err error) {
// 	if uid == "" || uid == "0" {
// 		return "", fmt.Errorf("invalid uid:%s", uid)
// 	}
// 	return uid, nil
// }

// func (tmpKnower) Uid(id string) (uid string, err error) {
// 	if id == "" || id == "0" {
// 		return "", fmt.Errorf("invalid id:%s", id)
// 	}
// 	return id, nil
// }

func (tk *tmpKnower) ServeAddr(m box.Message) (addr string, err error) {
	addr = string(m.GetAny().MessageName().Name())
	if addr == tk.invalidKey {
		addr, err = "", fmt.Errorf("invalid")
	}
	return
}

type tmpDialer struct {
	CP         chan<- *net.TmpConn
	invalidKey string
}

func (td *tmpDialer) Dial(ctx context.Context, addr string) (c net.Conn, err error) {
	var tc *net.TmpConn
	if addr == td.invalidKey {
		err = fmt.Errorf("invalid addr")
	} else {
		tc = &net.TmpConn{SendPP: make(chan box.Message), RecvPP: make(chan box.Message)}
	}
	td.CP <- tc
	return tc, err
}

func TestClientset(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cp, hc := make(chan *net.TmpConn, 1), make(chan net.TokenMessage, 16)
	cs := net.NewClientset(ctx, &tmpDialer{CP: cp, invalidKey: "S_Example"}, &tmpKnower{invalidKey: "G_Example"})
	cs.RouteMux().Handle("AnyTest1", net.FuncHandler(func(tm net.TokenMessage) {
		hc <- tm
	}))

	err := cs.AutoSend(net.TokenMessage{Token: one.TokenSet().Derive("25"), T: box.New(box.WithMessageValue(&testdata.AnyTest1{Name: "at1"}))})
	assert.Equal(t, err, nil)

	tc := <-cp
	m := <-tc.SendPP
	v, _ := m.GetAny().UnmarshalNew()
	assert.Equal(t, v.(*testdata.AnyTest1).GetName(), "at1")

	err = cs.AutoSend(net.TokenMessage{Token: one.TokenSet().Derive("25"), T: box.New(box.WithMessageValue(&testdata.AnyTest1{Name: "at11"}))})
	assert.Equal(t, err, nil)
	m = <-tc.SendPP
	v, _ = m.GetAny().UnmarshalNew()
	assert.Equal(t, v.(*testdata.AnyTest1).GetName(), "at11")

	t.Run("cannot_marshal_addr", func(t *testing.T) {
		err = cs.AutoSend(net.TokenMessage{Token: one.TokenSet().Derive("25"), T: box.New(box.WithMessageValue(&testdata.G_Example{Value: "gg"}))})
		assert.Equal(t, err.Error(), "invalid")

		//no connect
		// assert.Equal(t, <-cp, nil)
	})

	t.Run("cannot_connect", func(t *testing.T) {
		err = cs.AutoSend(net.TokenMessage{Token: one.TokenSet().Derive("25"), T: box.New(box.WithMessageValue(&testdata.S_Example{Value: "ss"}))})
		assert.Equal(t, err, nil, "can load client")

		tc := <-cp
		var emptytc *net.TmpConn
		assert.Equal(t, tc, emptytc, "connect failed")

		<-time.After(time.Millisecond * 10) //wait for remove from
		keys := make(chan string, 10)
		syncm := net.Export_clientsetm(cs)
		syncm.Range(func(key, value any) bool {
			keys <- key.(string)
			return true
		})
		assert.Equal(t, len(keys), 1, keys)
		key := <-keys
		assert.Equal(t, key, "AnyTest1")
	})

	// syncm := net.Export_clientsetm(cs)
	// syncm.Range()

	err = cs.AutoSend(net.TokenMessage{Token: one.TokenSet().Derive("26"), T: box.New(box.WithMessageValue(&testdata.P_Example{Value: "pp1"}))})
	assert.Equal(t, err, nil)
	cnt := 0
	syncm := net.Export_clientsetm(cs)
	syncm.Range(func(key, value any) bool {
		cnt++
		return true
	})
	assert.Equal(t, cnt, 2, "add P_Example client")

	tc2 := <-cp
	m = <-tc2.SendPP
	v, _ = m.GetAny().UnmarshalNew()
	assert.Equal(t, v.(*testdata.P_Example).GetValue(), "pp1")

	tc2.RecvPP <- box.New(box.WithMessageValue(&testdata.AnyTest1{Name: "at1r"}))
	tm := <-hc
	assert.Equal(t, tm.Token.Value(box.ContextkeyTokenUid).(string), "26")
	v, _ = tm.T.GetAny().UnmarshalNew()
	assert.Equal(t, v.(*testdata.AnyTest1).GetName(), "at1r")

	cs.Close()
	_, ok := <-tc.SendPP
	assert.Equal(t, ok, false, "conn closed after clientset closed")
	_, ok = <-tc2.SendPP
	assert.Equal(t, ok, false, "conn closed after clientset closed")
}

// func TestClient(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	cp := make(chan *net.TmpConn)
// 	cli := net.NewClient(ctx, &tmpDialer{CP: cp}, &tmpKnower{invalidKey: "G_Example"})
// 	cli.AutoSend(net.TokenMessage{Token: one.TokenSet().Derive("101"), T: box.New(box.WithMessageValue(&testdata.Test1{Name: "t1"}))})
// }

// type tmpRest struct{}

// func (tmpRest) Get(ctx context.Context, req *anypb.Any) (rsp *anypb.Any, err error) {
// 	if req == nil {
// 		err = fmt.Errorf("failed")
// 	} else {
// 		rsp, err = anypb.New(&testdata.AnyTest1{Name: "success"})
// 	}
// 	return
// }

// func (tmpRest) Post(ctx context.Context, req *anypb.Any) (out *emptypb.Empty, err error) {
// 	if req == nil {
// 		err = fmt.Errorf("failed")
// 	} else {
// 		out = &emptypb.Empty{}
// 	}
// 	return
// }

// func TestClientConn(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	t.Run("unimplement", func(t *testing.T) {
// 		cc := net.NewClientConn()
// 		assert.Equal(t, cc.Close(), net.ErrUnimplementApi)
// 		assert.Equal(t, cc.Send(nil), net.ErrUnimplementApi)
// 		_, err := cc.Recv()
// 		assert.Equal(t, err, net.ErrUnimplementApi)
// 		_, err = cc.Get(ctx, nil)
// 		assert.Equal(t, err, net.ErrUnimplementApi)
// 		_, err = cc.Post(ctx, nil)
// 		assert.Equal(t, err, net.ErrUnimplementApi)
// 	})
// 	t.Run("only conn", func(t *testing.T) {
// 		rpp, spp := make(chan *box.Message, 1), make(chan *box.Message, 1)
// 		cc := net.NewClientConn(net.WithConn(net.Export_newTmpConn(rpp, spp)))
// 		_, err := cc.Get(ctx, nil)
// 		assert.Equal(t, err, net.ErrUnimplementApi)
// 		_, err = cc.Post(ctx, nil)
// 		assert.Equal(t, err, net.ErrUnimplementApi)

// 		m, _ := box.New(&testdata.AnyTest1{Name: "a1"})
// 		err = cc.Send(m)
// 		assert.Equal(t, err, nil, err)
// 		m1 := <-spp
// 		v, _ := m1.GetAny().UnmarshalNew()
// 		assert.Equal(t, v.(*testdata.AnyTest1).GetName(), "a1")

// 		m, _ = box.New(&testdata.AnyTest2{Hope: "h2"})
// 		rpp <- m
// 		m2, err := cc.Recv()
// 		assert.Equal(t, err, nil, err)
// 		v, _ = m2.GetAny().UnmarshalNew()
// 		assert.Equal(t, v.(*testdata.AnyTest2).GetHope(), "h2")

// 		err = cc.Close()
// 		assert.Equal(t, err, nil, err)

// 		_, err = cc.Recv()
// 		assert.Equal(t, err, io.EOF)
// 	})
// 	t.Run("only Rest", func(t *testing.T) {
// 		cc := net.NewClientConn(net.WithRest(tmpRest{}))
// 		assert.Equal(t, cc.Close(), net.ErrUnimplementApi)
// 		assert.Equal(t, cc.Send(nil), net.ErrUnimplementApi)
// 		_, err := cc.Recv()
// 		assert.Equal(t, err, net.ErrUnimplementApi)
// 		_, err = cc.Get(ctx, new(anypb.Any))
// 		assert.Equal(t, err, nil)
// 		_, err = cc.Post(ctx, new(anypb.Any))
// 		assert.Equal(t, err, nil)
// 	})
// }
