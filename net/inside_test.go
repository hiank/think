package net

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/net/testdata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

type testConn struct {
	k string
	// identity uint64
	recvPP <-chan pb.M
	sendPP chan<- pb.M
}

// func (tc *testConn) GetIdentity() uint64 {
// 	return tc.identity
// }

func (tc *testConn) Recv() (m pb.M, err error) {
	m, ok := <-tc.recvPP
	if !ok {
		err = io.EOF
	}
	return
}

func (tc *testConn) Send(m pb.M) error {
	tc.sendPP <- m
	return nil
}

func (tc *testConn) Close() error {
	if tc.sendPP != nil {
		close(tc.sendPP)
		tc.sendPP = nil
	}
	return nil
}

// func TestAnyCoder(t *testing.T) {
// 	var ac coder.AnyBytes
// 	d, err := ac.Encode(&testdata.AnyTest1{Name: "test1"})
// 	assert.Assert(t, err == nil, err)

// 	v, err := ac.Decode(d)
// 	assert.Assert(t, err == nil, err)
// 	assert.Equal(t, v.(*testdata.AnyTest1).Name, "test1")

// 	amsg, err := anypb.New(&testdata.AnyTest1{Name: "hiank"})
// 	assert.Assert(t, err == nil, err)
// 	d, err = ac.Encode(amsg)
// 	assert.Assert(t, err == nil, err)
// 	v, _ = ac.Decode(d)
// 	assert.Equal(t, v.(*testdata.AnyTest1).Name, "hiank")

// 	_, err = ac.Encode(&testConn{})
// 	assert.Assert(t, err != nil)

// 	_, err = ac.Decode([]byte{1, 2})
// 	assert.Assert(t, err != nil)
// }

func TestConnpool(t *testing.T) {
	t.Run("lookErr", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cp, terr := &connpool{ctx: ctx}, fmt.Errorf("terr")
		assert.Equal(t, cp.lookErr(terr), terr)
		assert.Equal(t, cp.lookErr(nil), error(nil))

		cancel()
		assert.Equal(t, cp.lookErr(nil), ctx.Err())
		assert.Equal(t, cp.lookErr(terr), ctx.Err())
	})

	t.Run("AddConn", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cp := newConnpool(ctx, nil) //&connpool{ctx: ctx}

		pp := make(chan pb.M)
		c1 := &testConn{sendPP: pp, k: "c1"}
		cp.AddConn("1", c1)
		_, ok := cp.m.Load("1")
		assert.Assert(t, ok)
		// assert.Equal(t, mv.(Conn).(*testConn), c1)
		wait, pp2 := make(chan byte), make(chan pb.M)
		go func() {
			cp.AddConn("1", &testConn{sendPP: pp2, k: "c2"})
			close(wait)
		}()
		_, ok = <-pp
		assert.Assert(t, !ok, "closed by Close")

		<-wait

		cnt := 0
		cp.m.Range(func(key, value interface{}) bool {
			cnt++
			return true
		})
		assert.Equal(t, cnt, 1)

		mv, _ := cp.m.Load("1")
		k := mv.(*fatconn).conn.(*testConn).k
		assert.Equal(t, k, "c2")

		mv.(Conn).Close()

		cnt = 0
		cp.m.Range(func(key, value interface{}) bool {
			cnt++
			return true
		})
		assert.Equal(t, cnt, 0)
	})

	t.Run("loopCheck", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cp := &connpool{ctx: ctx}
		go func(t *testing.T) {
			// cp.loopCheck()
			cancel()
		}(t)

		cp.loopCheck()
	})
	t.Run("Send", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cp := newConnpool(ctx, nil) //&connpool{ctx: ctx}

		pp, pp2 := make(chan pb.M), make(chan pb.M)
		cp.AddConn("1", &testConn{sendPP: pp, k: "c1"})
		cp.AddConn("2", &testConn{sendPP: pp2, k: "c2"})

		cnt := 0
		cp.m.Range(func(key, value interface{}) bool {
			cnt++
			return true
		})
		assert.Equal(t, cnt, 2)

		err := cp.Send(&testdata.AnyTest1{Name: "hiank"})
		assert.Assert(t, err == nil, err)

		v := <-pp
		v2 := <-pp2
		amsg, _ := anypb.New(&testdata.AnyTest1{Name: "hiank"})
		tv, err := proto.Marshal(amsg)
		assert.Assert(t, err == nil, err)
		assert.Equal(t, len(v.Bytes()), len(v2.Bytes()))
		assert.Equal(t, len(v.Bytes()), len(tv))

		for i, b := range tv {
			assert.Equal(t, b, v.Bytes()[i])
			assert.Equal(t, b, v2.Bytes()[i])
		}

		err = cp.Send(&testdata.AnyTest2{Hope: "hope"}, "3")
		assert.Equal(t, err, ErrNoConn)

		tids := []string{"1", "2"}
		err = cp.Send(&testdata.AnyTest2{Hope: "hope"}, tids...)
		assert.Equal(t, err, nil)
	})
}

func TestCopy(t *testing.T) {
	l := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	copy(l[5:], l[6:])
	l = l[:len(l)-1]
	// t.Log(l)
	for i, v := range []int{0, 1, 2, 3, 4, 6, 7, 8, 9} {
		assert.Equal(t, l[i], v)
	}
}

func TestContext(t *testing.T) {
	ctx1, cancel1 := context.WithCancel(context.Background())
	_, cancel2 := context.WithCancel(ctx1)

	cancel2()
	cancel2() //repeated calls are ok

	_, cancel2 = context.WithCancel(ctx1)
	cancel1()
	assert.Assert(t, ctx1.Err() != nil)
	cancel2() //call cancel after parent context closed is ok
}

func TestPrivateServer(t *testing.T) {
	// newServer(nil, nil)
	t.Run("lookErr", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		srv := &server{ctx: ctx, cancel: cancel, connpool: newConnpool(ctx, nil)}
		err := srv.lookErr(nil)
		assert.Equal(t, err, nil)
		err = srv.lookErr(fmt.Errorf("ws-err"))
		assert.Equal(t, err.Error(), "ws-err")
		cancel()
		err = srv.lookErr(nil)
		assert.Equal(t, err, context.Canceled)
		err = srv.lookErr(fmt.Errorf("ig-err"))
		assert.Equal(t, err, context.Canceled)
	})
}

// 	t.Run("handleConn", func(t *testing.T) {
// 		ctx, cancel := context.WithCancel(context.Background())
// 		defer cancel()
// 		recvPP, sendPP := make(chan *anypb.Any), make(chan *anypb.Any)
// 		srv := &server{ctx: ctx, cancel: cancel}
// 		note := make(chan bool)
// 		go func() {
// 			srv.handleConn(&testConn{identity: 1, recvPP: recvPP, sendPP: sendPP})
// 			note <- true
// 		}()

// 		select {
// 		case <-note:
// 			assert.Assert(t, false, "should block here")
// 		case <-time.After(time.Millisecond * 100):
// 			assert.Assert(t, true, "should block in handleConn")
// 		}

// 		recvPP2, sendPP2 := make(chan *anypb.Any), make(chan *anypb.Any)
// 		srv.handleConn(&testConn{identity: 1, recvPP: recvPP2, sendPP: sendPP2})
// 		assert.Assert(t, true, "same identity, would not block")

// 		go func() {
// 			srv.handleConn(&testConn{identity: 2, recvPP: recvPP2, sendPP: sendPP2})
// 		}()
// 		<-time.After(time.Millisecond * 100) //NOTE: wait identity 2 ready

// 		// assert.Equal(t, )
// 		val, _ := srv.m.Load(uint64(2))
// 		assert.Equal(t, val.(Conn).GetIdentity(), uint64(2))

// 		val, _ = srv.m.Load(uint64(1))
// 		assert.Equal(t, val.(Conn).GetIdentity(), uint64(1))
// 	})

// }
