package net

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/hiank/think/net/box"
	"github.com/hiank/think/run"
	"google.golang.org/protobuf/proto"
	"k8s.io/klog/v2"
)

const (
	ErrNonTargetConn     = run.Err("net: non target conn")
	ErrNonTargetIdentity = run.Err("net: non target identity for send")
	ErrClosed            = run.Err("net: closed")
)

type connset struct {
	// ctx context.Context
	h Handler  //Handler for revc message
	m sync.Map //conn map
}

func newConnset(h Handler) *connset {
	return &connset{h: h}
}

//loadOrStore new taskConn
func (cp *connset) loadOrStore(ctx context.Context, id string, connect Connect) (lc *liteConn, err error) {
	v, loaded := cp.m.LoadOrStore(id, &liteConn{})
	if lc = v.(*liteConn); !loaded {
		err = initialize(ctx, lc, connect, cp.loopRecv, func() { cp.m.Delete(id) })
	}
	return
}

//lookErr check error and print
func lookErr(err error) error {
	if err != nil && err != io.EOF {
		klog.Warning(err)
	}
	return err
}

//loopRecv loop read from given conn
func (cs *connset) loopRecv(receiver Receiver, tk box.Token) {
	for {
		d, err := receiver.Recv()
		if err = lookErr(err); err != nil {
			return
		}
		go cs.h.Route(TokenMessage{T: d, Token: tk.Fork()})
	}
}

//broadcast send message to all conn
func (cs *connset) broadcast(m box.Message) (err error) {
	cs.m.Range(func(_, value any) bool {
		if tmperr := lookErr(value.(Conn).Send(m)); tmperr != nil {
			klog.Warning("net: connset Send error:", tmperr)
			err = tmperr
		}
		return true
	})
	return
}

//multiSend send message to multi conn
func (cs *connset) multiSend(m box.Message, tis ...string) (err error) {
	if len(tis) == 0 {
		return ErrNonTargetIdentity
	}
	km := make(map[any]byte)
	for _, k := range tis {
		km[k] = 1
	}
	cs.m.Range(func(key, value any) bool {
		if _, ok := km[key]; ok {
			if tmperr := lookErr(value.(Conn).Send(m)); tmperr != nil {
				klog.Warning("net: connset Send error:", tmperr)
				err = tmperr
			}
			delete(km, key)
		}
		return len(km) > 0
	})
	for k := range km {
		klog.Warning("net: cannot found target connect:", k)
		err = ErrNonTargetConn
	}
	return
}

//close clear conn store (close all conn)
func (cs *connset) close() {
	cs.m.Range(func(_, value any) bool {
		lookErr(value.(Conn).Close())
		return true
	})
}

type Connect func(context.Context) (TokenConn, error)

//initialize liteConn
func initialize(ctx context.Context, lc *liteConn, connect Connect, loopRecv func(Receiver, box.Token), doneHook func()) error {
	if ctx.Err() != nil {
		doneHook()
		return ctx.Err()
	}
	ctx, cancel := context.WithCancel(ctx)
	lc.tasker = run.NewTasker(ctx, time.Second*10)
	lc.closer = run.NewOnceCloser(func() error {
		doneHook()
		cancel()
		return nil
	})
	return lc.tasker.Add(run.NewLiteTask(func(lc *liteConn) (err error) {
		tc, err := connect(ctx)
		if err == nil {
			if err = run.FrontErr(tc.Token.Err, ctx.Err, func() error {
				return lc.ready(tc, doneHook)
			}); err == nil {
				go func() {
					loopRecv(tc.T, tc.Token)
					lc.Close()
				}()
				return
			}
			tc.T.Close()
		}
		klog.Warningln("net: connect failed", err)
		return lc.Close()
	}, lc))
}

//liteConn lightweight Conn
//contians a tasker. execute connect->send... in sequence
//NOTE: unsafe to call Recv() from multiple goroutine. because Receiver maybe reset
type liteConn struct {
	tasker    run.Tasker
	tc        TokenConn
	closer    io.Closer
	onceReset sync.Once //reset closer once (after connect success)
}

func (lc *liteConn) Send(m box.Message) error {
	return lc.tasker.Add(run.NewLiteTask(func(m box.Message) (err error) {
		if err = lc.tc.T.Send(m); err != nil {
			if err != io.EOF {
				klog.Warningf("conn write error: %v", err)
			}
			//close conn and tell connset to delete this conn
			lc.Close()
		}
		return
	}, m))
}

//Recv unimplemented. start loopRecv in initialize when connect success
func (lc *liteConn) Recv() (box.Message, error) {
	return nil, ErrUnimplementedApi
}

//
func (lc *liteConn) Close() error {
	lc.onceReset.Do(func() {
		//avoid closer reset after Close executed
	})
	return lc.closer.Close()
}

func (lc *liteConn) ready(tc TokenConn, doneHook func()) (err error) {
	err = ErrClosed
	lc.onceReset.Do(func() {
		healthy := run.NewHealthy()
		lc.tc, lc.closer = tc, run.NewHealthyCloser(healthy, func() { tc.Token.Close() })
		go healthy.Monitoring(tc.Token, func() {
			tc.T.Close()
			doneHook()
		})
		err = nil
	})
	return
}

type RouteMux struct {
	m sync.Map
}

//Handle register Handler for k
//k must be string/proto.Message value
func (rm *RouteMux) Handle(k any, h Handler) {
	var sk string
	switch v := k.(type) {
	case string:
		sk = v
	case proto.Message:
		sk = string(v.ProtoReflect().Descriptor().FullName())
	default:
		klog.Warning("net: unsupport k value type")
		return
	}
	rm.m.Store(sk, h)
}

func (rm *RouteMux) Route(tt TokenMessage) {
	k := string(tt.T.GetAny().MessageName().Name())
	mv, loaded := rm.m.Load(k)
	if !loaded {
		if mv, loaded = rm.m.Load(DefaultHandler); !loaded {
			klog.Warning("cannot find handler for handle message recv by conn: ", k)
			return
		}
	}
	mv.(Handler).Route(tt)
}
