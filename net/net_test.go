package net_test

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/hiank/think/net"
	"github.com/hiank/think/pbtest"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

func TestMessage(t *testing.T) {
	//
	t.Run("panic if not messagevalue/messagebytes given", func(t *testing.T) {
		defer func(t *testing.T) {
			r := recover()
			assert.Equal(t, r, net.ErrNonMessageValue)
		}(t)
		net.NewMessage(net.WithMessageToken(nil))
	})

	msg := net.NewMessage(net.WithMessageToken(nil), net.WithMessageBytes(nil))
	assert.DeepEqual(t, msg.Any().GetValue(), []byte(nil))

	msg = net.NewMessage(net.WithMessageToken(nil), net.WithMessageBytes(nil), net.WithMessageValue(&pbtest.AnyTest1{Name: "at1"}))
	assert.DeepEqual(t, msg.Any().GetValue(), []byte(nil))

	msg = net.NewMessage(net.WithMessageToken(nil), net.WithMessageValue(&pbtest.AnyTest1{Name: "at2"}))
	assert.Equal(t, string(msg.Any().MessageName()), "AnyTest1")
	pm, _ := msg.Any().UnmarshalNew()
	assert.Equal(t, pm.(*pbtest.AnyTest1).GetName(), "at2")

	t.Run("panic nil message value", func(t *testing.T) {
		defer func(t *testing.T) {
			r := recover()
			assert.Assert(t, r != nil)
		}(t)
		net.NewMessage(net.WithMessageValue(nil))
	})
	t.Run("panic invalid bytes", func(t *testing.T) {
		defer func(t *testing.T) {
			r := recover()
			assert.Assert(t, r != nil)
		}(t)
		net.NewMessage(net.WithMessageBytes([]byte("invalid")))
	})

	b, _ := proto.Marshal(&pbtest.AnyTest1{Name: "at3"})
	msg = net.NewMessage(net.WithMessageBytes(b))
	// assert.Equal(t, err, nil, "非anypb.Any bytes 解码不会报错，但无结果")
	assert.DeepEqual(t, msg.Any().GetValue(), []byte(nil))
}

func TestAnyNew(t *testing.T) {
	amsg, _ := anypb.New(&pbtest.AnyTest1{Name: "cw"})
	v2, err := anypb.New(amsg)
	assert.Equal(t, err, nil, err)
	v3, err := v2.UnmarshalNew()
	assert.Equal(t, err, nil, err)
	amsg, ok := v3.(*anypb.Any)
	assert.Assert(t, ok)
	v4, _ := amsg.UnmarshalNew()
	assert.Equal(t, v4.(*pbtest.AnyTest1).GetName(), "cw")

	var amsg2 anypb.Any = *amsg
	v5, _ := (&amsg2).UnmarshalNew()
	assert.Equal(t, v5.(*pbtest.AnyTest1).GetName(), "cw")
}

func TestRouteMux(t *testing.T) {
	rm := &net.RouteMux{}
	nt := make(chan *net.Message, 1)
	rm.Handle("S_Example", net.FuncHandler(func(tt *net.Message) {
		nt <- tt
	}))

	err := rm.Handle(1, net.FuncHandler(func(tt *net.Message) {}))
	assert.Equal(t, err, net.ErrUnsupportValueType)
	//get a warning log here. unsupport handle type

	m := net.NewMessage(net.WithMessageValue(&pbtest.S_Example{Value: "route"}))
	rm.Route(m)

	tt := <-nt
	v, _ := tt.Any().UnmarshalNew()
	assert.Equal(t, v.(*pbtest.S_Example).GetValue(), "route")

	m = net.NewMessage(net.WithMessageValue(&pbtest.G_Example{Value: "gg"}))
	rm.Route(m)
	//get a warning log here

	dnt := make(chan *net.Message, 1)
	rm.Handle(net.DefaultHandler, net.FuncHandler(func(tt *net.Message) {
		dnt <- tt
	}))

	m = net.NewMessage(net.WithMessageValue(&pbtest.G_Example{Value: "gg2"}))
	rm.Route(m)
	tt = <-dnt
	v, _ = tt.Any().UnmarshalNew()
	assert.Equal(t, v.(*pbtest.G_Example).GetValue(), "gg2", "handle by default handler")
}

type tmpKnower struct {
	invalidKey string
}

func (tk *tmpKnower) ServeAddr(m *net.Message) (addr string, err error) {
	addr = string(m.Any().MessageName().Name())
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
		tc = &net.TmpConn{SendPP: make(chan *net.Message), RecvPP: make(chan proto.Message), Tk: net.Tokenset.Derive("default")}
	}
	td.CP <- tc
	return tc, err
}

func TestClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cp, hc := make(chan *net.TmpConn, 1), make(chan *net.Message, 16)
	cs := net.NewClient(ctx, &tmpDialer{CP: cp, invalidKey: "S_Example"}, &tmpKnower{invalidKey: "G_Example"})
	cs.RouteMux().Handle("AnyTest1", net.FuncHandler(func(tm *net.Message) {
		hc <- tm
	}))

	m := net.NewMessage(net.WithMessageValue(&pbtest.AnyTest1{Name: "at1"}), net.WithMessageToken(net.Tokenset.Derive("25")))
	err := cs.AutoSend(m)
	assert.Equal(t, err, nil)

	tc := <-cp
	// return
	m = <-tc.SendPP
	v, _ := m.Any().UnmarshalNew()
	assert.Equal(t, v.(*pbtest.AnyTest1).GetName(), "at1")

	m = net.NewMessage(net.WithMessageValue(&pbtest.AnyTest1{Name: "at11"}), net.WithMessageToken(net.Tokenset.Derive("25")))
	err = cs.AutoSend(m)
	assert.Equal(t, err, nil)
	m = <-tc.SendPP
	v, _ = m.Any().UnmarshalNew()
	assert.Equal(t, v.(*pbtest.AnyTest1).GetName(), "at11")

	t.Run("cannot_marshal_addr", func(t *testing.T) {
		m := net.NewMessage(net.WithMessageValue(&pbtest.G_Example{Value: "gg"}), net.WithMessageToken(net.Tokenset.Derive("25")))
		err = cs.AutoSend(m)
		assert.Equal(t, err.Error(), "invalid")

		//no connect
		// assert.Equal(t, <-cp, nil)
	})

	t.Run("cannot_connect", func(t *testing.T) {
		m := net.NewMessage(net.WithMessageValue(&pbtest.S_Example{Value: "ss"}), net.WithMessageToken(net.Tokenset.Derive("25")))
		err = cs.AutoSend(m)
		assert.Equal(t, err, nil, "can load client")

		tc := <-cp
		var emptytc *net.TmpConn
		assert.Equal(t, tc, emptytc, "connect failed")

		<-time.After(time.Millisecond * 10) //wait for remove from
		keys := make(chan string, 10)
		syncm := net.Export_clientm(cs)
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

	m = net.NewMessage(net.WithMessageValue(&pbtest.P_Example{Value: "pp1"}), net.WithMessageToken(net.Tokenset.Derive("26")))
	err = cs.AutoSend(m)
	assert.Equal(t, err, nil)
	cnt := 0
	syncm := net.Export_clientm(cs)
	syncm.Range(func(key, value any) bool {
		cnt++
		return true
	})
	assert.Equal(t, cnt, 2, "add P_Example client")

	tc2 := <-cp
	m = <-tc2.SendPP
	v, _ = m.Any().UnmarshalNew()
	assert.Equal(t, v.(*pbtest.P_Example).GetValue(), "pp1")

	tc2.RecvPP <- &pbtest.AnyTest1{Name: "at1r"}
	tm := <-hc
	assert.Equal(t, tm.Token().ToString(), "default")
	v, _ = tm.Any().UnmarshalNew()
	assert.Equal(t, v.(*pbtest.AnyTest1).GetName(), "at1r")

	cs.Close()
	_, ok := <-tc.SendPP
	assert.Equal(t, ok, false, "conn closed after clientset closed")
	_, ok = <-tc2.SendPP
	assert.Equal(t, ok, false, "conn closed after clientset closed")
}

type tmpListener struct {
	// once chan uint64
	connPP chan net.Conn
}

func (tl *tmpListener) Accept() (iac net.Conn, err error) {
	iac, ok := <-tl.connPP
	if !ok {
		err = io.EOF
	}
	return
}

func (tl *tmpListener) Close() error {
	close(tl.connPP)
	return nil
}

func TestServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cpp, router := make(chan net.Conn, 10), &net.RouteMux{}
	lis := &tmpListener{connPP: cpp}
	srv := net.NewServer(ctx, lis, router)
	vpp := make(chan *net.Message, 20)
	router.Handle(&pbtest.G_Example{}, net.FuncHandler(func(tt *net.Message) {
		vpp <- tt
	}))

	wait := make(chan bool)
	go func() {
		srv.ListenAndServe()
		close(wait)
	}()

	recvPP, sendPP := make(chan proto.Message), make(chan *net.Message)
	cpp <- &net.TmpConn{RecvPP: recvPP, SendPP: sendPP, Tk: net.Tokenset.Derive("110")}
	recvPP1, sendPP1 := make(chan proto.Message), make(chan *net.Message)
	cpp <- &net.TmpConn{RecvPP: recvPP1, SendPP: sendPP1, Tk: net.Tokenset.Derive("111")}

	t.Run("accept-recv", func(t *testing.T) {
		recvPP <- &pbtest.G_Example{Value: "g1"}
		recvPP <- &pbtest.S_Example{Value: "s1"}

		<-time.After(time.Millisecond * 10)
		assert.Equal(t, len(vpp), 1) //only G_Example could be response
		tt := <-vpp
		assert.Equal(t, tt.Token().ToString(), "110")
		v, _ := tt.Any().UnmarshalNew()
		assert.Equal(t, v.(*pbtest.G_Example).GetValue(), "g1")
	})
	t.Run("send", func(t *testing.T) {
		srv.Send(&pbtest.Test1{Name: "t1"})
		m1, m2 := <-sendPP, <-sendPP1
		assert.Equal(t, m1, m2)
		v, _ := m1.Any().UnmarshalNew()
		assert.Equal(t, v.(*pbtest.Test1).GetName(), "t1")

		err := srv.Send(&pbtest.Test2{Hope: "t2"}, "110", "113")
		assert.Assert(t, err != nil)

		m1 = <-sendPP
		v, _ = m1.Any().UnmarshalNew()
		assert.Equal(t, v.(*pbtest.Test2).GetHope(), "t2")

		select {
		case <-sendPP1:
			assert.Assert(t, false, "do not send to the id")
		case <-time.After(time.Millisecond * 100):
		}
	})

	srv.Close()
	_, err := lis.Accept()
	assert.Equal(t, err, io.EOF, "listener closed")

	err = srv.Send(&pbtest.AnyTest1{Name: "at1"})
	assert.Equal(t, err, context.Canceled, "closed")
}

func TestRedefineOutValue(t *testing.T) {
	func1 := func() (int, int) {
		return 1, 11
	}
	func2 := func(t *testing.T) (outVal int) {
		val, outVal := func1()
		assert.Equal(t, val, 1)
		assert.Equal(t, outVal, 11)
		outVal = 12
		if val == 1 {
			val1, outVal := func1()
			assert.Equal(t, val1, 1)
			assert.Equal(t, outVal, 11)
		}
		return
	}
	assert.Equal(t, func2(t), 12, "if代码块中 outVal作用域只在其中")
}

func TestDeleteNonItem(t *testing.T) {
	m := map[int]int{1: 11}
	delete(m, 2)
	assert.Equal(t, len(m), 1, "delete not existed key would not panic")
}

type testKey string

type testKey1 string

func TestInterfaceType(t *testing.T) {
	var key testKey = "key"
	var it any = key
	val, ok := it.(testKey1)
	assert.Assert(t, !ok)
	assert.Equal(t, val, testKey1(""))
}

func TestForcedConversion(t *testing.T) {
	data := &pbtest.AnyTest1{Name: "hp"}
	var tmp any = data
	val, _ := tmp.(*pbtest.AnyTest2)
	// assert.Assert(t, !ok)
	assert.Assert(t, val == nil)

	var arr []int
	arr2 := make([]int, 0)
	arr2 = append(arr2, arr...)
	assert.Equal(t, len(arr2), 0)
	assert.Assert(t, arr == nil)
}
