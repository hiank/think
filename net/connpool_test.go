package net

import (
	"context"
	"fmt"
	"io"
	"time"

	"testing"

	"github.com/hiank/think/net/box"
	"github.com/hiank/think/net/testdata"
	"gotest.tools/v3/assert"
)

type tmpConn struct {
	k string
	// identity uint64
	recvPP chan *box.Message
	sendPP chan<- *box.Message
}

func (tc *tmpConn) Recv() (m *box.Message, err error) {
	m, ok := <-tc.recvPP
	if !ok {
		err = io.EOF
	}
	return
}

func (tc *tmpConn) Send(m *box.Message) error {
	tc.sendPP <- m
	return nil
}

func (tc *tmpConn) Close() error {
	if tc.sendPP != nil {
		close(tc.sendPP)
		tc.sendPP = nil
	}
	if tc.recvPP != nil {
		close(tc.recvPP)
	}
	return nil
}

func TestTaskConn(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	t.Run("canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		pp1, rm := make(chan *box.Message), make(chan bool)
		// _, err := newTaskConn(ctx, func(ctx context.Context) (Conn, error) {
		// 	return &tmpConn{sendPP: pp1}, nil
		// }, func() { close(rm) })
		tc := &taskConn{}
		tc.init(ctx)
		tc.monitorConn(&tmpConn{sendPP: pp1}, func() { close(rm) })
		<-rm
		// assert.Assert(t, err != nil)
	})

	t.Run("dial fialed", func(t *testing.T) {
		rm := make(chan bool)
		tc := &taskConn{}
		tc.init(ctx)
		tc.monitorConnect(func(ctx context.Context) (Conn, error) {
			return nil, fmt.Errorf("invalid")
		}, func() { close(rm) })
		<-rm

		ctx, cancel := context.WithCancel(ctx)
		// cancel()
		sp, rm := make(chan *box.Message), make(chan bool)
		tc = &taskConn{}
		tc.init(ctx)
		tc.monitorConnect(func(ctx context.Context) (Conn, error) {
			<-time.After(time.Millisecond * 100)
			return &tmpConn{sendPP: sp}, nil
		}, func() { close(rm) })

		<-time.After(time.Millisecond * 50)
		cancel()
		_, ok := <-sp
		assert.Equal(t, ok, false)
		_, ok = <-rm
		assert.Equal(t, ok, false)
	})

	pp1, rm := make(chan *box.Message), make(chan bool)
	rc := make(chan *box.Message)
	tc := &taskConn{}
	tc.init(ctx)
	tc.monitorConnect(func(ctx context.Context) (Conn, error) {
		<-time.After(time.Millisecond * 100)
		return &tmpConn{sendPP: pp1, recvPP: rc}, nil
	}, func() { close(rm) })

	m, _ := box.New(&testdata.AnyTest1{Name: "a1"})
	err := tc.Send(m)
	assert.Equal(t, err, nil, err)
	t1 := time.Now().UnixMilli()
	m1 := <-pp1
	tt := time.Now().UnixMilli() - t1
	assert.Assert(t, tt > 99, "wait until dial completed", tt)
	v, _ := m1.GetAny().UnmarshalNew()
	assert.Equal(t, v.(*testdata.AnyTest1).GetName(), "a1")

	go func() {
		m, _ := box.New(&testdata.AnyTest2{Hope: "h1"})
		rc <- m
	}()
	m2, err := tc.Recv()
	assert.Equal(t, err, nil, err)
	v2, _ := m2.GetAny().UnmarshalNew()
	assert.Equal(t, v2.(*testdata.AnyTest2).GetHope(), "h1")
}

type tmpHandleValue struct {
	id string
	m  *box.Message
}

func TestConnpool(t *testing.T) {
	t.Run("lookErr", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cp, terr := newConnpool(ctx, nil), fmt.Errorf("terr")
		assert.Equal(t, cp.lookErr(terr), terr)
		assert.Equal(t, cp.lookErr(nil), error(nil))

		cancel()
		assert.Equal(t, cp.lookErr(nil), ctx.Err())
		assert.Equal(t, cp.lookErr(terr), ctx.Err())
	})

	t.Run("add", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cp := newConnpool(ctx, nil) //&connpool{ctx: ctx}

		pp := make(chan *box.Message)
		c1 := &tmpConn{sendPP: pp, k: "c1"}
		cp.add("1", c1)
		_, ok := cp.m.Load("1") //.Load("1")
		assert.Assert(t, ok)
		wait, pp2 := make(chan byte), make(chan *box.Message)
		go func(t *testing.T) {
			<-time.After(time.Millisecond * 100)
			cp.add("1", &tmpConn{sendPP: pp2, k: "c2"})
			close(wait)
		}(t)
		_, ok = <-pp
		assert.Assert(t, !ok, "closed by Close")

		<-wait

		cnt := 0
		cp.m.Range(func(key, value any) bool {
			cnt++
			return true
		})
		assert.Equal(t, cnt, 1)

		mv, _ := cp.m.Load("1")
		// c := net.Export_taskConnValueC(mv)
		// k := mv.(*tmpConn).k
		c := mv.(*taskConn).c
		assert.Equal(t, c.(*tmpConn).k, "c2")

		mv.(Conn).Close()
		// <-time.After(time.Millisecond * 10)

		cnt = 0
		cp.m.Range(func(key, value any) bool {
			cnt++
			return true
		})
		assert.Equal(t, cnt, 0)

	})

	t.Run("broadcast-multiSend", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cp := newConnpool(ctx, nil) //&connpool{ctx: ctx}

		pps := []chan *box.Message{
			make(chan *box.Message, 1),
			make(chan *box.Message, 1),
			make(chan *box.Message, 1),
		}
		cs := make([]*tmpConn, len(pps))

		for i, pp := range pps {
			cs[i] = &tmpConn{sendPP: pp, k: fmt.Sprintf("c%d", i+1)}
			func(idx int) {
				cp.add(fmt.Sprintf("%d", idx+1), cs[idx])
			}(i)
		}

		m, _ := box.New(&testdata.AnyTest1{Name: "bc"})
		err := cp.broadcast(m)
		assert.Equal(t, err, nil)

		for _, pp := range pps {
			v, ok := <-pp
			assert.Equal(t, ok, true)
			// assert.DeepEqual(t, v, m)
			msg, _ := v.GetAny().UnmarshalNew()
			assert.Equal(t, msg.(*testdata.AnyTest1).GetName(), "bc")
		}

		m, _ = box.New(&testdata.AnyTest2{Hope: "hp"})
		err = cp.multiSend(m, "1", "3")
		assert.Equal(t, err, nil)

		v := <-pps[0]
		msg, _ := v.GetAny().UnmarshalNew()
		assert.Equal(t, msg.(*testdata.AnyTest2).GetHope(), "hp")

		v = <-pps[2]
		msg, _ = v.GetAny().UnmarshalNew()
		assert.Equal(t, msg.(*testdata.AnyTest2).GetHope(), "hp")

		err = cp.multiSend(m, "4")
		assert.Equal(t, err, ErrNonTargetConn)

		err = cp.multiSend(m)
		assert.Equal(t, err, ErrNonTargetIdentity)
	})

	t.Run("loopRecv", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		router := &RouteMux{}
		cp := newConnpool(ctx, router) //&connpool{ctx: ctx}

		hc := make(chan tmpHandleValue)
		router.Handle("AnyTest1", HandlerFunc(func(id string, m *box.Message) {
			hc <- tmpHandleValue{id: id, m: m}
		}))

		router.Handle(new(testdata.AnyTest2), HandlerFunc(func(id string, m *box.Message) {
			hc <- tmpHandleValue{id: id, m: m}
		}))

		pps := []chan *box.Message{
			make(chan *box.Message, 1),
			make(chan *box.Message, 1),
			make(chan *box.Message, 1),
		}
		cs := make([]*tmpConn, len(pps))

		for i, pp := range pps {
			cs[i] = &tmpConn{recvPP: pp, k: fmt.Sprintf("c%d", i+1)}
			func(idx int) {
				cp.add(fmt.Sprintf("%d", idx+1), cs[idx])
			}(i)
		}

		m, _ := box.New(&testdata.AnyTest1{Name: "bc"})
		pps[0] <- m

		tv1 := <-hc
		assert.Equal(t, tv1.id, "1")
		v1, _ := tv1.m.GetAny().UnmarshalNew()
		assert.Equal(t, v1.(*testdata.AnyTest1).GetName(), "bc")

		m, _ = box.New(&testdata.Test1{Name: "t1"})
		pps[0] <- m

		select {
		case tv1 = <-hc:
			assert.Assert(t, false, "non handler for Test1")
		case <-time.After(time.Millisecond * 10):
		}

		m, _ = box.New(&testdata.AnyTest2{Hope: "hp"})
		pps[2] <- m

		tv1 = <-hc
		assert.Equal(t, tv1.id, "3")
		v1, _ = tv1.m.GetAny().UnmarshalNew()
		assert.Equal(t, v1.(*testdata.AnyTest2).GetHope(), "hp")

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
