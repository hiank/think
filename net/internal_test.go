package net

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"testing"

	"github.com/hiank/think/net/box"
	"github.com/hiank/think/net/one"
	"github.com/hiank/think/net/testdata"
	"gotest.tools/v3/assert"
)

type TmpConn struct {
	// k string
	// identity uint64
	RecvPP chan box.Message
	SendPP chan box.Message
}

func (tc *TmpConn) Recv() (m box.Message, err error) {
	m, ok := <-tc.RecvPP
	if !ok {
		err = io.EOF
	}
	return
}

func (tc *TmpConn) Send(m box.Message) error {
	tc.SendPP <- m
	return nil
}

func (tc *TmpConn) Close() error {
	if tc.SendPP != nil {
		close(tc.SendPP)
		// tc.SendPP = nil
	}
	if tc.RecvPP != nil {
		close(tc.RecvPP)
	}
	return nil
}

func tmpConnect(ctx context.Context) (tc TokenConn, err error) {
	// tc.Token, _ = one.TokenSet().Build()
	tc.Token = box.NewToken(ctx)
	tc.T = &TmpConn{RecvPP: make(chan box.Message), SendPP: make(chan box.Message)}
	return
}

func TestLiteConnInitialize(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("ctx canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		c := make(chan bool)
		err := initialize(ctx, &liteConn{}, tmpConnect, func(r Receiver, t box.Token) {}, func() { close(c) })
		<-c
		assert.Assert(t, err != nil, "context canceled")
	})

	t.Run("non-initialize", func(t *testing.T) {
		lc := &liteConn{}
		defer func(t *testing.T) {
			r := recover()
			assert.Assert(t, r != nil)
		}(t)
		lc.Send(box.New(box.WithMessageValue(&testdata.AnyTest1{})))
	})

	t.Run("close before connect complete", func(t *testing.T) {
		lc, tc, wait := &liteConn{}, &TmpConn{SendPP: make(chan box.Message)}, make(chan bool)
		err := initialize(ctx, lc, func(ctx context.Context) (TokenConn, error) {
			close(wait)
			<-time.After(time.Millisecond * 10)
			return TokenConn{T: tc, Token: one.TokenSet().Derive("empty")}, nil
		}, func(r Receiver, t box.Token) {}, func() {})
		assert.Equal(t, err, nil, nil)

		lc.Send(box.New(box.WithMessageValue(&testdata.AnyTest1{Name: "at1"})))
		<-wait //wait until connect task start
		lc.Close()
		_, ok := <-tc.SendPP
		assert.Assert(t, !ok, "closed")
	})

	t.Run("delay send", func(t *testing.T) {
		lc, tc, pp := &liteConn{}, &TmpConn{SendPP: make(chan box.Message)}, make(chan int, 3)
		err := initialize(ctx, lc, func(ctx context.Context) (TokenConn, error) {
			pp <- 1
			<-time.After(time.Millisecond * 10)
			pp <- 2
			return TokenConn{T: tc, Token: one.TokenSet().Derive("empty")}, nil
		}, func(r Receiver, t box.Token) {}, func() {})
		assert.Equal(t, err, nil, nil)

		err = lc.Send(box.New(box.WithMessageValue(&testdata.AnyTest1{Name: "at1"})))
		assert.Equal(t, err, nil, err)
		<-tc.SendPP
		assert.Equal(t, len(pp), 2)
	})

	t.Run("connect failed", func(t *testing.T) {
		lc, pp, hook := &liteConn{}, make(chan int, 3), make(chan bool)
		err := initialize(ctx, lc, func(ctx context.Context) (TokenConn, error) {
			<-pp
			<-time.After(time.Millisecond * 10)
			return TokenConn{}, fmt.Errorf("connect failed")
		}, func(r Receiver, t box.Token) {}, func() { close(hook) })
		assert.Equal(t, err, nil, nil)

		err = lc.Send(box.New(box.WithMessageValue(&testdata.AnyTest1{Name: "at1"})))
		assert.Equal(t, err, nil, err)
		pp <- 1 //notice to do connect

		<-hook //auto closed after connect failed
	})

}

var makeConnect = func(id string, tc Conn) Connect {
	return func(ctx context.Context) (TokenConn, error) {
		<-time.After(time.Millisecond * 10)
		return TokenConn{T: tc, Token: one.TokenSet().Derive(id)}, nil
	}
}

func TestConnset(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("load-close", func(t *testing.T) {
		// ctx, cancel := context.WithCancel(ctx)
		// defer cancel()
		router := &RouteMux{}
		cs := newConnset(router)
		connected, s1 := make(chan int, 3), make(chan box.Message)
		lc, _ := cs.loadOrStore(ctx, "110", func(ctx context.Context) (TokenConn, error) {
			connected <- 1
			<-time.After(time.Millisecond * 10)
			return TokenConn{T: &TmpConn{SendPP: s1}, Token: one.TokenSet().Derive("110")}, nil
		})
		s2 := make(chan box.Message)
		lc2, _ := cs.loadOrStore(ctx, "110", func(ctx context.Context) (TokenConn, error) {
			connected <- 2
			<-time.After(time.Millisecond * 10)
			return TokenConn{T: &TmpConn{SendPP: s2}, Token: one.TokenSet().Derive("110")}, nil
		})
		assert.Equal(t, lc, lc2, "")

		s3 := make(chan box.Message)
		lc3, _ := cs.loadOrStore(ctx, "111", func(ctx context.Context) (TokenConn, error) {
			connected <- 3
			<-time.After(time.Millisecond * 10)
			return TokenConn{T: &TmpConn{SendPP: s3}, Token: one.TokenSet().Derive("111")}, nil
		})
		assert.Assert(t, lc3 != lc)

		<-time.After(time.Millisecond * 10)
		assert.Equal(t, len(connected), 2)
		assert.Equal(t, <-connected+<-connected, 4)

		cs.close()
		// cancel()
		<-s1
		<-s3 //all conn closed
	})

	router := &RouteMux{}
	cs := newConnset(router)
	defer cs.close()
	tcs := []*TmpConn{
		{SendPP: make(chan box.Message), RecvPP: make(chan box.Message)},
		{SendPP: make(chan box.Message), RecvPP: make(chan box.Message)},
		{SendPP: make(chan box.Message), RecvPP: make(chan box.Message)},
	}
	cs.loadOrStore(ctx, "110", makeConnect("110", tcs[0]))
	cs.loadOrStore(ctx, "112", makeConnect("112", tcs[1]))
	cs.loadOrStore(ctx, "111", makeConnect("111", tcs[2]))

	t.Run("broadcast-multiSend", func(t *testing.T) {
		err := cs.broadcast(box.New(box.WithMessageValue(&testdata.AnyTest1{Name: "at1"})))
		assert.Equal(t, err, nil)

		for _, tc := range tcs {
			m := <-tc.SendPP
			v, _ := m.GetAny().UnmarshalNew()
			assert.Equal(t, v.(*testdata.AnyTest1).GetName(), "at1")
		}

		err = cs.multiSend(box.New(box.WithMessageValue(&testdata.AnyTest2{Hope: "h1"})), "110", "112", "113")
		assert.Equal(t, err, ErrNonTargetConn)
		for _, tc := range tcs[:2] {
			m := <-tc.SendPP
			v, _ := m.GetAny().UnmarshalNew()
			assert.Equal(t, v.(*testdata.AnyTest2).GetHope(), "h1")
		}
		select {
		case <-tcs[2].SendPP:
			assert.Assert(t, false, "not target conn")
		case <-time.After(time.Millisecond * 100):
			assert.Assert(t, true, "")
		}
	})

	cnt := 10
	for i := 0; i < cnt; i++ {
		go func(str string) {
			cs.multiSend(box.New(box.WithMessageValue(&testdata.Test2{Hope: str})), "112")
		}(strconv.Itoa(i))
	}

	want := 0
L:
	for {
		select {
		case m := <-tcs[1].SendPP:
			v, _ := m.GetAny().UnmarshalNew()
			i, _ := strconv.ParseInt(v.(*testdata.Test2).GetHope(), 10, 32)
			want |= (1 << int(i))
		case <-time.After(time.Millisecond * 100):
			break L
		}
	}
	assert.Equal(t, want, (1<<cnt)-1)

	cache := make(chan TokenMessage, 20)
	// want = 0
	router.Handle(&testdata.Test1{}, FuncHandler(func(tt TokenMessage) {
		// assert.Equal(t, tt.Token.Value(box.ContextkeyTokenUid).(string), "113")
		// v, _ := tt.T.GetAny().UnmarshalNew()
		cache <- tt
	}))

	for i := 0; i < cnt; i++ {
		go func(str string) {
			tcs[2].RecvPP <- box.New(box.WithMessageValue(&testdata.Test1{Name: str}))
		}(strconv.Itoa(i))
	}
	want = 0
	i := 0
	for tt := range cache {
		assert.Equal(t, tt.Token.Value(box.ContextkeyTokenUid).(string), "111")
		v, _ := tt.T.GetAny().UnmarshalNew()
		iv, _ := strconv.Atoi(v.(*testdata.Test1).GetName())
		want |= (1 << iv)
		i++
		if i == 10 {
			break
		}
	}
	assert.Equal(t, want, (1<<cnt)-1)

	<-time.After(time.Millisecond * 100)
	assert.Equal(t, len(cache), 0, "no more message recv")
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
