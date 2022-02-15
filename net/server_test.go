package net_test

import (
	"context"
	"io"
	"testing"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/testdata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

type testHandler struct {
	out chan *net.Doc
}

func (tch *testHandler) Route(id string, d *net.Doc) {
	tch.out <- d
	// return nil
}

type testListener struct {
	// once chan uint64
	connPP chan *net.IAC
}

func (tl *testListener) Accept() (iac net.IAC, err error) {
	tc, ok := <-tl.connPP
	if ok {
		iac = *tc
	} else {
		err = io.EOF
	}
	return
}

func (tl *testListener) Close() error {
	close(tl.connPP)
	return nil
}

func TestNewServer(t *testing.T) {
	// ctx, cancel := context.WithCancel(context.TODO())
	// defer cancel()
	router := &net.RouteMux{}
	srv := net.NewServer(&testListener{connPP: make(chan *net.IAC)}, router)
	go func() {
		srv.Close()
	}()
	err := srv.ListenAndServe()
	assert.Equal(t, err, io.EOF)

	err = srv.Close()
	assert.Equal(t, err, context.Canceled, "close could called repeatedly")
}

func TestDoc(t *testing.T) {
	doc, err := net.MakeDoc(&testdata.G_Example{})
	assert.Assert(t, err == nil, err)
	assert.Equal(t, doc.TypeName(), "G_Example", doc.TypeName())

	_, err = net.MakeDoc(11)
	assert.Assert(t, err != nil, "param for makedoc should be a proto.Message")

	// b := doc.Bytes()
	var amsg anypb.Any
	err = proto.Unmarshal(doc.Bytes(), &amsg)
	assert.Assert(t, err == nil, err)
}

func TestRouteMux(t *testing.T) {
	rm := &net.RouteMux{}
	nt := make(chan *net.Doc, 1)
	rm.Handle("S_Example", net.HandlerFunc(func(s string, d *net.Doc) {
		nt <- d
	}))

	d, _ := net.MakeDoc(&testdata.S_Example{Value: "route"})
	rm.Route("tmp", d)

	d = <-nt
	v, _ := d.Any().UnmarshalNew()
	assert.Equal(t, v.(*testdata.S_Example).GetValue(), "route")
}

func TestServer(t *testing.T) {

	handlerPP := make(chan *net.Doc)
	accept, router := make(chan *net.IAC), &net.RouteMux{}
	router.Handle("", &testHandler{out: handlerPP})
	srv := net.NewServer(&testListener{connPP: accept}, router)
	defer srv.Close()
	go srv.ListenAndServe()

	t.Run("accept-recv-send", func(t *testing.T) {
		recvPP, sendPP := make(chan *net.Doc), make(chan *net.Doc)
		tc := net.Export_newTestConn(recvPP, sendPP) //&testConn{recvPP: recvPP, sendPP: sendPP, identity: 1}
		accept <- &net.IAC{ID: "1", Conn: tc}

		msg := &testdata.S_Example{Value: "pp"}
		any, _ := anypb.New(msg)
		d, _ := proto.Marshal(any)
		ndoc, _ := net.MakeDoc(d)
		recvPP <- ndoc

		doc := <-handlerPP
		// assert.Equal(t, carrier.GetIdentity(), uint64(1))
		val, _ := doc.Any().UnmarshalNew()
		assert.Equal(t, val.(*testdata.S_Example).GetValue(), "pp")

		err := srv.Send(&testdata.AnyTest1{Name: "test1"}, "2")
		assert.Assert(t, err != nil, "no id 2 conn dialed")

		// any, _ = anypb.New(&testdata.AnyTest1{Name: "ws"})
		go func(t *testing.T) {
			err = srv.Send(&testdata.AnyTest1{Name: "ws"}, "1")
			assert.Equal(t, err, nil)
		}(t)
		doc = <-sendPP
		// any = <-sendPP
		proto.Unmarshal(doc.Bytes(), any)
		val, _ = any.UnmarshalNew()
		assert.Equal(t, val.(*testdata.AnyTest1).GetName(), "ws")
	})
	// t.Run("NewServer")
}

func TestForcedConversion(t *testing.T) {
	data := &testdata.AnyTest1{Name: "hp"}
	var tmp interface{} = data
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
	var it interface{} = key
	val, ok := it.(testKey1)
	assert.Assert(t, !ok)
	assert.Equal(t, val, testKey1(""))
}
