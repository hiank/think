package net

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/hiank/think/net/box"
	"github.com/hiank/think/run"
	"k8s.io/klog/v2"
)

const (
	ErrNonTargetConn     = run.Err("net: non target conn")
	ErrNonTargetIdentity = run.Err("net: non target identity for send")
	ErrInvalidDocParam   = run.Err("net: invalid doc param: should be bytes/proto.Message")
)

type connpool struct {
	ctx context.Context
	h   Handler  //Handler for revc message
	m   sync.Map //conn map
	io.Closer
}

func newConnpool(ctx context.Context, h Handler) (cp *connpool) {
	ctx, cancel := context.WithCancel(ctx)
	healthy := run.NewHealthy()
	cp = &connpool{
		ctx:    ctx,
		h:      h,
		Closer: run.NewHealthyCloser(healthy, cancel),
	}
	go healthy.Monitoring(ctx, func() {
		cp.m.Range(func(_, value any) bool {
			cp.lookErr(value.(Conn).Close())
			return true
		})
		///this connpool will invalid after here
	})
	return
}

//add new taskConn
func (cp *connpool) add(id string, conn Conn) {
	tc := &taskConn{}
	tc.init(cp.ctx)
	tc.monitorConn(conn, func() {
		cp.m.Delete(id)
	})
	for v, loaded := cp.m.LoadOrStore(id, tc); loaded; v, loaded = cp.m.LoadOrStore(id, tc) {
		//wait until stored conn deleted
		v.(io.Closer).Close()
	}
	go cp.loopRecv(id, tc)
}

//loadOrStore new taskConn
func (cp *connpool) loadOrStore(id string, connect Connect) (c *taskConn, err error) {
	v, loaded := cp.m.LoadOrStore(id, &taskConn{})
	if c = v.(*taskConn); !loaded {
		if err = cp.ctx.Err(); err == nil {
			c.init(cp.ctx)
			c.monitorConnect(connect, func() {
				cp.m.Delete(id)
			})
			go cp.loopRecv(id, c)
		}
	}
	return
}

//lookErr check error and print
func (cp *connpool) lookErr(err error) error {
	if cp.ctx.Err() != nil {
		err = cp.ctx.Err()
	}
	if err != nil && err != io.EOF {
		klog.Warning(err)
	}
	return err
}

//loopRecv loop read from given conn
func (cp *connpool) loopRecv(id string, conn Conn) {
	defer conn.Close()
	for {
		d, err := conn.Recv()
		if err = cp.lookErr(err); err != nil {
			return
		}
		go cp.h.Route(id, d)
	}
}

//broadcast send message to all conn
func (cp *connpool) broadcast(m *box.Message) (err error) {
	cp.m.Range(func(_, value any) bool {
		if tmperr := cp.lookErr(value.(Conn).Send(m)); tmperr != nil {
			klog.Warning("net: connpool Send error:", tmperr)
			err = tmperr
		}
		return true
	})
	return
}

//multiSend send message to multi conn
func (cp *connpool) multiSend(m *box.Message, tis ...string) (err error) {
	if len(tis) == 0 {
		return ErrNonTargetIdentity
	}
	km := make(map[any]byte)
	for _, k := range tis {
		km[k] = 1
	}
	cp.m.Range(func(key, value any) bool {
		if _, ok := km[key]; ok {
			if tmperr := cp.lookErr(value.(Conn).Send(m)); tmperr != nil {
				klog.Warning("net: connpool Send error:", tmperr)
				err = tmperr
			}
			delete(km, key)
		}
		return len(km) > 0
	})
	for _, k := range km {
		klog.Warning("net: cannot found target connect:", k)
		err = ErrNonTargetConn
	}
	return
}

// type recover struct {

// }

type recoverFunc func() (*box.Message, error)

func (rf recoverFunc) Recv() (*box.Message, error) {
	return rf()
}

// type Dial func(context.Context) (Conn, error)

type Connect func(context.Context) (Conn, error)

type taskConn struct {
	ctx context.Context
	// connectx context.Context
	c       Conn
	tasker  run.Tasker
	healthy *run.Healthy
	Receiver
	io.Closer
}

func (tc *taskConn) init(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	healthy := run.NewHealthy()
	*tc = taskConn{
		ctx:     ctx,
		tasker:  run.NewTasker(ctx, time.Second*5),
		healthy: healthy,
		Closer:  run.NewHealthyCloser(healthy, cancel),
	}
}

func (tc *taskConn) monitorConn(c Conn, doneHook func()) {
	tc.c, tc.Receiver = c, recoverFunc(func() (*box.Message, error) {
		return c.Recv()
	})
	go tc.healthy.Monitoring(tc.ctx, func() {
		tc.c.Close()
		doneHook()
	})
}

func (tc *taskConn) monitorConnect(connect Connect, doneHook func()) {
	connectx, cancel := context.WithCancel(tc.ctx)
	tc.Receiver = recoverFunc(func() (out *box.Message, err error) {
		<-connectx.Done()
		//check weather taskConn closed
		if err = tc.ctx.Err(); err == nil {
			out, err = tc.c.Recv()
		}
		return
	})
	go tc.healthy.Monitoring(connectx, func() {
		defer doneHook()
		if tc.c != nil {
			<-tc.ctx.Done()
			tc.c.Close()
		}
	})
	tc.tasker.Add(run.NewLiteTask(func(tc *taskConn) (err error) {
		c, err := connect(connectx)
		if err != nil {
			klog.Warningln("net: taskConn dial failed", err)
			tc.Close()
		} else if tc.ctx.Err() == nil { //not closed outside
			tc.c = c
			cancel()
		} else { //closed outside
			c.Close()
		}
		return
	}, tc)) //first task is connect
}

//Send add a send task
func (tc *taskConn) Send(m *box.Message) (err error) {
	return tc.tasker.Add(run.NewLiteTask(tc.taskSend, m))
}

//taskSend for Task
func (tc *taskConn) taskSend(m *box.Message) (err error) {
	if werr := tc.c.Send(m); werr != nil {
		if werr != io.EOF {
			klog.Warningf("conn write error: %s", werr.Error())
		}
		//close conn and tell connpool to delete this conn
		tc.Close()
	}
	return
}
