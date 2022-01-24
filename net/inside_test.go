package net

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

type testConn struct {
	identity uint64
	recvPP   <-chan *anypb.Any
	sendPP   chan<- *anypb.Any
}

func (tc *testConn) GetIdentity() uint64 {
	return tc.identity
}

func (tc *testConn) Recv() (any *anypb.Any, err error) {
	any, ok := <-tc.recvPP
	if !ok {
		err = io.EOF
	}
	return
}

func (tc *testConn) Send(any *anypb.Any) error {
	tc.sendPP <- any
	return nil
}

func (tc *testConn) Close() error {
	return nil
}

func TestPrivateServer(t *testing.T) {
	// newServer(nil, nil)
	t.Run("lookErr", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		srv := &server{ctx: ctx, cancel: cancel}
		err := srv.lookErr(nil)
		assert.Equal(t, err, nil)
		err = srv.lookErr(errors.New("ws-err"))
		assert.Equal(t, err.Error(), "ws-err")
		cancel()
		err = srv.lookErr(nil)
		assert.Equal(t, err, context.Canceled)
		err = srv.lookErr(errors.New("ig-err"))
		assert.Equal(t, err, context.Canceled)
	})

	t.Run("handleConn", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		recvPP, sendPP := make(chan *anypb.Any), make(chan *anypb.Any)
		srv := &server{ctx: ctx, cancel: cancel}
		note := make(chan bool)
		go func() {
			srv.handleConn(&testConn{identity: 1, recvPP: recvPP, sendPP: sendPP})
			note <- true
		}()

		select {
		case <-note:
			assert.Assert(t, false, "should block here")
		case <-time.After(time.Millisecond * 100):
			assert.Assert(t, true, "should block in handleConn")
		}

		recvPP2, sendPP2 := make(chan *anypb.Any), make(chan *anypb.Any)
		srv.handleConn(&testConn{identity: 1, recvPP: recvPP2, sendPP: sendPP2})
		assert.Assert(t, true, "same identity, would not block")

		go func() {
			srv.handleConn(&testConn{identity: 2, recvPP: recvPP2, sendPP: sendPP2})
		}()
		<-time.After(time.Millisecond * 100) //NOTE: wait identity 2 ready

		// assert.Equal(t, )
		val, _ := srv.m.Load(uint64(2))
		assert.Equal(t, val.(Conn).GetIdentity(), uint64(2))

		val, _ = srv.m.Load(uint64(1))
		assert.Equal(t, val.(Conn).GetIdentity(), uint64(1))
	})

}
