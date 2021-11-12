package net_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/hiank/think/net"
	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/net/testdata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

type testMessageHandler struct {
	out chan proto.Message
}

func (tch *testMessageHandler) Handle(id uint64, msg proto.Message) {
	tch.out <- msg
}

type testCarrierHandler struct {
	out chan *pb.Carrier
}

// func (tch *tes)

func (tch *testCarrierHandler) Handle(carrier *pb.Carrier) {
	tch.out <- carrier
	// return nil
}

type testListener struct {
	// once chan uint64
	connPP chan net.IConn
}

func (tl *testListener) Accept() (conn net.IConn, err error) {
	tc, ok := <-tl.connPP
	if ok {
		conn = tc //&testConn{identity: identity, recvPP: make(chan *anypb.Any), sendPP: make(chan *anypb.Any)}
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
	srv := net.NewServer(&testListener{connPP: make(chan net.IConn)}, &testCarrierHandler{})
	go func() {
		srv.Close()
	}()
	err := srv.ListenAndServe()
	assert.Equal(t, err, context.Canceled)

	err = srv.Close()
	assert.Equal(t, err, context.Canceled, "close could called repeatedly")
}

func TestServer(t *testing.T) {
	// srv := net.NewServer(nil)
	// err := srv.ListenAndServe()
	accept, handlerPP := make(chan net.IConn), make(chan *pb.Carrier)
	srv := net.NewServer(&testListener{connPP: accept}, &testCarrierHandler{handlerPP})
	defer srv.Close()
	go srv.ListenAndServe()

	t.Run("accept-recv-send", func(t *testing.T) {
		recvPP, sendPP := make(chan *anypb.Any), make(chan *anypb.Any)
		tc := net.Export_newTestConn(1, recvPP, sendPP) //&testConn{recvPP: recvPP, sendPP: sendPP, identity: 1}
		accept <- tc

		msg := &testdata.S_Example{Value: "pp"}
		any, _ := anypb.New(msg)
		recvPP <- any

		carrier := <-handlerPP
		assert.Equal(t, carrier.GetIdentity(), uint64(1))
		val, _ := carrier.GetMessage().UnmarshalNew()
		assert.Equal(t, val.(*testdata.S_Example).GetValue(), "pp")

		err := srv.Send(&pb.Carrier{Identity: 2})
		assert.Assert(t, err != nil)

		any, _ = anypb.New(&testdata.AnyTest1{Name: "ws"})
		go func(t *testing.T) {
			err = srv.Send(&pb.Carrier{Identity: 1, Message: any})
			assert.Equal(t, err, nil)
		}(t)
		any = <-sendPP
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

func TestServerWithHandleMux(t *testing.T) {
	t.Run("non-option", func(t *testing.T) {
		accept, handlerPP, handleMux := make(chan net.IConn), make(chan proto.Message), net.NewHandleMux()
		srv := net.NewServer(&testListener{connPP: accept}, handleMux)
		defer srv.Close()
		go srv.ListenAndServe()

		handleMux.Look("AnyTest1", &testMessageHandler{handlerPP})

		recvPP, sendPP := make(chan *anypb.Any), make(chan *anypb.Any)
		tc := net.Export_newTestConn(1, recvPP, sendPP) //&testConn{recvPP: recvPP, sendPP: sendPP, identity: 1}
		accept <- tc

		msg := &testdata.AnyTest2{Hope: "pp"}
		any, _ := anypb.New(msg)
		recvPP <- any
		select {
		case <-handlerPP:
			assert.Assert(t, false, "no handler for the message")
		case <-time.After(time.Millisecond * 100):
			// default:
			assert.Assert(t, true)
		}

		msg2 := &testdata.AnyTest1{Name: "ts1"}
		any, _ = anypb.New(msg2)
		recvPP <- any
		val := <-handlerPP
		// val.GetIdentity()
		// val, _ := carrier.GetMessage().UnmarshalNew()
		assert.Equal(t, val.(*testdata.AnyTest1).GetName(), "ts1")
	})
	t.Run("WithDefaultHandler", func(t *testing.T) {
		accept, handlerPP := make(chan net.IConn), make(chan *pb.Carrier)
		handleMux := net.NewHandleMux(net.WithDefaultHandler(&testCarrierHandler{handlerPP}))
		srv := net.NewServer(&testListener{connPP: accept}, handleMux)
		defer srv.Close()
		go srv.ListenAndServe()

		handlerPP2 := make(chan proto.Message)
		// handleMux.Look("AnyTest1", &testMessageHandler{handlerPP2})
		handleMux.LookObject(new(testdata.AnyTest1), &testMessageHandler{handlerPP2})

		recvPP, sendPP := make(chan *anypb.Any), make(chan *anypb.Any)
		tc := net.Export_newTestConn(1, recvPP, sendPP) //&testConn{recvPP: recvPP, sendPP: sendPP, identity: 1}
		accept <- tc

		msg := &testdata.AnyTest2{Hope: "pp"}
		any, _ := anypb.New(msg)
		recvPP <- any
		carrier := <-handlerPP
		val, _ := carrier.GetMessage().UnmarshalNew()
		assert.Equal(t, val.(*testdata.AnyTest2).GetHope(), "pp")

		msg2 := &testdata.AnyTest1{Name: "ts2"}
		any, _ = anypb.New(msg2)
		recvPP <- any
		val = <-handlerPP2
		// val, _ = carrier.GetMessage().UnmarshalNew()
		assert.Equal(t, val.(*testdata.AnyTest1).GetName(), "ts2")
	})
	t.Run("WithConverter", func(t *testing.T) {
		accept := make(chan net.IConn)
		handleMux := net.NewHandleMux(net.WithConverter(net.FuncCarrierConverter(func(c *pb.Carrier) (string, bool) {
			return "anyTest2", true
		})))
		srv := net.NewServer(&testListener{connPP: accept}, handleMux)
		defer srv.Close()
		go srv.ListenAndServe()

		handlerPP1, handlerPP2 := make(chan proto.Message), make(chan proto.Message)
		handleMux.Look("AnyTest1", &testMessageHandler{handlerPP1})
		handleMux.Look("anyTest2", &testMessageHandler{handlerPP2})

		recvPP, sendPP := make(chan *anypb.Any), make(chan *anypb.Any)
		tc := net.Export_newTestConn(1, recvPP, sendPP) //&testConn{recvPP: recvPP, sendPP: sendPP, identity: 1}
		accept <- tc

		msg := &testdata.AnyTest2{Hope: "pp"}
		any, _ := anypb.New(msg)
		recvPP <- any
		select {
		case <-handlerPP1:
			assert.Assert(t, false, "no handler for the message")
		case <-handlerPP2:
			assert.Assert(t, true, "key is anyTest2")
		case <-time.After(time.Millisecond * 100):
			// default:
			assert.Assert(t, false)
		}

		msg2 := &testdata.AnyTest1{Name: "ts1"}
		any, _ = anypb.New(msg2)
		recvPP <- any
		val := <-handlerPP2
		// val.GetIdentity()
		// val, _ := carrier.GetMessage().UnmarshalNew()
		assert.Equal(t, val.(*testdata.AnyTest1).GetName(), "ts1")
	})
}
