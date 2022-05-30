package net_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/box"
	"github.com/hiank/think/net/one"
	"github.com/hiank/think/net/testdata"
	"github.com/hiank/think/run"
	"gotest.tools/v3/assert"
)

// type tmpHandler struct {
// 	out chan net.TokenMessage
// }

// func (tch *tmpHandler) Route(tt net.TokenMessage) {
// 	tch.out <- tt
// 	// return nil
// }

type tmpListener struct {
	// once chan uint64
	connPP chan net.TokenConn
}

func (tl *tmpListener) Accept() (iac net.TokenConn, err error) {
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

func TestNewServer(t *testing.T) {

	router := &net.RouteMux{}
	srv := net.NewServer(&tmpListener{connPP: make(chan net.TokenConn)}, router)
	go func() {
		srv.Close()
	}()
	err := srv.ListenAndServe()
	assert.Equal(t, err, io.EOF)

	err = srv.Close()
	assert.Equal(t, err, run.ErrBeenClosed, "")
}

func TestRouteMux(t *testing.T) {
	rm := &net.RouteMux{}
	nt := make(chan net.TokenMessage, 1)
	rm.Handle("S_Example", net.FuncHandler(func(tt net.TokenMessage) {
		nt <- tt
	}))

	rm.Handle(1, net.FuncHandler(func(tt net.TokenMessage) {}))
	//get a warning log here. unsupport handle type

	d := box.New(box.WithMessageValue(&testdata.S_Example{Value: "route"}))
	// rm.Route("tmp", d)
	rm.Route(net.TokenMessage{T: d})

	tt := <-nt
	v, _ := tt.T.GetAny().UnmarshalNew()
	assert.Equal(t, v.(*testdata.S_Example).GetValue(), "route")

	rm.Route(net.TokenMessage{T: box.New(box.WithMessageValue(&testdata.G_Example{Value: "gg"}))})
	//get a warning log here

	dnt := make(chan net.TokenMessage, 1)
	rm.Handle("", net.FuncHandler(func(tt net.TokenMessage) {
		dnt <- tt
	}))

	rm.Route(net.TokenMessage{T: box.New(box.WithMessageValue(&testdata.G_Example{Value: "gg2"}))})
	tt = <-dnt
	v, _ = tt.T.GetAny().UnmarshalNew()
	assert.Equal(t, v.(*testdata.G_Example).GetValue(), "gg2", "handle by default handler")
}

func TestServer(t *testing.T) {
	cpp, router := make(chan net.TokenConn, 10), &net.RouteMux{}
	lis := &tmpListener{connPP: cpp}
	srv := net.NewServer(lis, router)
	vpp := make(chan net.TokenMessage, 20)
	router.Handle(&testdata.G_Example{}, net.FuncHandler(func(tt net.TokenMessage) {
		vpp <- tt
	}))

	wait := make(chan bool)
	go func() {
		srv.ListenAndServe()
		close(wait)
	}()

	recvPP, sendPP := make(chan box.Message), make(chan box.Message)
	cpp <- net.TokenConn{T: &net.TmpConn{RecvPP: recvPP, SendPP: sendPP}, Token: one.TokenSet().Derive("110")}
	recvPP1, sendPP1 := make(chan box.Message), make(chan box.Message)
	cpp <- net.TokenConn{T: &net.TmpConn{RecvPP: recvPP1, SendPP: sendPP1}, Token: one.TokenSet().Derive("111")}

	t.Run("accept-recv", func(t *testing.T) {
		recvPP <- box.New(box.WithMessageValue(&testdata.G_Example{Value: "g1"}))
		recvPP <- box.New(box.WithMessageValue(&testdata.S_Example{Value: "s1"}))

		<-time.After(time.Millisecond * 10)
		assert.Equal(t, len(vpp), 1) //only G_Example could be response
		tt := <-vpp
		assert.Equal(t, tt.Token.Value(box.ContextkeyTokenUid).(string), "110")
		v, _ := tt.T.GetAny().UnmarshalNew()
		assert.Equal(t, v.(*testdata.G_Example).GetValue(), "g1")
	})
	t.Run("send", func(t *testing.T) {
		srv.Send(&testdata.Test1{Name: "t1"})
		m1, m2 := <-sendPP, <-sendPP1
		assert.Equal(t, m1, m2)
		v, _ := m1.GetAny().UnmarshalNew()
		assert.Equal(t, v.(*testdata.Test1).GetName(), "t1")

		err := srv.Send(&testdata.Test2{Hope: "t2"}, "110", "113")
		assert.Equal(t, err, net.ErrNonTargetConn)

		m1 = <-sendPP
		v, _ = m1.GetAny().UnmarshalNew()
		assert.Equal(t, v.(*testdata.Test2).GetHope(), "t2")

		select {
		case <-sendPP1:
			assert.Assert(t, false, "do not send to the id")
		case <-time.After(time.Millisecond * 100):
		}
	})

	srv.Close()
	_, err := lis.Accept()
	assert.Equal(t, err, io.EOF, "listener closed")

	err = srv.Send(&testdata.AnyTest1{Name: "at1"})
	assert.Equal(t, err, context.Canceled, "closed")
}

func TestForcedConversion(t *testing.T) {
	data := &testdata.AnyTest1{Name: "hp"}
	var tmp any = data
	val, _ := tmp.(*testdata.AnyTest2)
	// assert.Assert(t, !ok)
	assert.Assert(t, val == nil)

	var arr []int
	arr2 := make([]int, 0)
	arr2 = append(arr2, arr...)
	assert.Equal(t, len(arr2), 0)
	assert.Assert(t, arr == nil)
}

// // func TestServerWithHandleMux(t *testing.T) {
// // 	t.Run("non-option", func(t *testing.T) {
// // 		accept, handlerPP, handleMux := make(chan *net.IAC), make(chan proto.Message), net.NewHandleMux()
// // 		srv := net.NewServer(&testListener{connPP: accept}, handleMux)
// // 		defer srv.Close()
// // 		go srv.ListenAndServe()

// // 		handleMux.Look("AnyTest1", &testMessageHandler{handlerPP})

// // 		recvPP, sendPP := make(chan *anypb.Any), make(chan *anypb.Any)
// // 		tc := net.Export_newTestConn(1, recvPP, sendPP) //&testConn{recvPP: recvPP, sendPP: sendPP, identity: 1}
// // 		accept <- tc

// // 		msg := &testdata.AnyTest2{Hope: "pp"}
// // 		any, _ := anypb.New(msg)
// // 		recvPP <- any
// // 		select {
// // 		case <-handlerPP:
// // 			assert.Assert(t, false, "no handler for the message")
// // 		case <-time.After(time.Millisecond * 100):
// // 			// default:
// // 			assert.Assert(t, true)
// // 		}

// // 		msg2 := &testdata.AnyTest1{Name: "ts1"}
// // 		any, _ = anypb.New(msg2)
// // 		recvPP <- any
// // 		val := <-handlerPP
// // 		// val.GetIdentity()
// // 		// val, _ := carrier.GetMessage().UnmarshalNew()
// // 		assert.Equal(t, val.(*testdata.AnyTest1).GetName(), "ts1")
// // 	})
// // 	t.Run("WithDefaultHandler", func(t *testing.T) {
// // 		accept, handlerPP := make(chan net.Conn), make(chan *pb.Carrier)
// // 		handleMux := net.NewHandleMux(net.WithDefaultHandler(&testCarrierHandler{handlerPP}))
// // 		srv := net.NewServer(&testListener{connPP: accept}, handleMux)
// // 		defer srv.Close()
// // 		go srv.ListenAndServe()

// // 		handlerPP2 := make(chan proto.Message)
// // 		// handleMux.Look("AnyTest1", &testMessageHandler{handlerPP2})
// // 		handleMux.LookObject(new(testdata.AnyTest1), &testMessageHandler{handlerPP2})

// // 		recvPP, sendPP := make(chan *anypb.Any), make(chan *anypb.Any)
// // 		tc := net.Export_newTestConn(1, recvPP, sendPP) //&testConn{recvPP: recvPP, sendPP: sendPP, identity: 1}
// // 		accept <- tc

// // 		msg := &testdata.AnyTest2{Hope: "pp"}
// // 		any, _ := anypb.New(msg)
// // 		recvPP <- any
// // 		carrier := <-handlerPP
// // 		val, _ := carrier.GetMessage().UnmarshalNew()
// // 		assert.Equal(t, val.(*testdata.AnyTest2).GetHope(), "pp")

// // 		msg2 := &testdata.AnyTest1{Name: "ts2"}
// // 		any, _ = anypb.New(msg2)
// // 		recvPP <- any
// // 		val = <-handlerPP2
// // 		// val, _ = carrier.GetMessage().UnmarshalNew()
// // 		assert.Equal(t, val.(*testdata.AnyTest1).GetName(), "ts2")
// // 	})
// // 	t.Run("WithConverter", func(t *testing.T) {
// // 		accept := make(chan net.Conn)
// // 		handleMux := net.NewHandleMux(net.WithConverter(net.FuncCarrierConverter(func(c *pb.Carrier) (string, bool) {
// // 			return "anyTest2", true
// // 		})))
// // 		srv := net.NewServer(&testListener{connPP: accept}, handleMux)
// // 		defer srv.Close()
// // 		go srv.ListenAndServe()

// // 		handlerPP1, handlerPP2 := make(chan proto.Message), make(chan proto.Message)
// // 		handleMux.Look("AnyTest1", &testMessageHandler{handlerPP1})
// // 		handleMux.Look("anyTest2", &testMessageHandler{handlerPP2})

// // 		recvPP, sendPP := make(chan *anypb.Any), make(chan *anypb.Any)
// // 		tc := net.Export_newTestConn(1, recvPP, sendPP) //&testConn{recvPP: recvPP, sendPP: sendPP, identity: 1}
// // 		accept <- tc

// // 		msg := &testdata.AnyTest2{Hope: "pp"}
// // 		any, _ := anypb.New(msg)
// // 		recvPP <- any
// // 		select {
// // 		case <-handlerPP1:
// // 			assert.Assert(t, false, "no handler for the message")
// // 		case <-handlerPP2:
// // 			assert.Assert(t, true, "key is anyTest2")
// // 		case <-time.After(time.Millisecond * 100):
// // 			// default:
// // 			assert.Assert(t, false)
// // 		}

// // 		msg2 := &testdata.AnyTest1{Name: "ts1"}
// // 		any, _ = anypb.New(msg2)
// // 		recvPP <- any
// // 		val := <-handlerPP2
// // 		// val.GetIdentity()
// // 		// val, _ := carrier.GetMessage().UnmarshalNew()
// // 		assert.Equal(t, val.(*testdata.AnyTest1).GetName(), "ts1")
// // 	})
// // }

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
	// var key1 testKey1 = "key1"
	var key testKey = "key"
	var it any = key
	val, ok := it.(testKey1)
	assert.Assert(t, !ok)
	assert.Equal(t, val, testKey1(""))
}
