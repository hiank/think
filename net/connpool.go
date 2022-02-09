package net

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/hiank/think/run"
	"k8s.io/klog/v2"
)

type connpool struct {
	ctx context.Context
	h   Handler     //Handler for revc message
	m   sync.Map    //conn map
	rm  chan string //for remove conn
	io.Closer
}

func newConnpool(ctx context.Context, h Handler) (cp *connpool) {
	ctx, cancel := context.WithCancel(ctx)
	cp = &connpool{
		ctx: ctx,
		h:   h,
		rm:  make(chan string),
	}
	cp.Closer = run.NewOnceCloser(func() error {
		close(cp.rm)
		cancel()
		return ctx.Err()
	})
	go cp.loopCheck()
	return
}

//AddConn add new conn to pool
//when existed same id, delete it and close it's conn before
//wait unit add success
func (cp *connpool) AddConn(id string, c Conn) {
	fc := newFatconn(cp.ctx, id, c, cp.rm)
	for {
		//wait until stored conn deleted
		v, loaded := cp.m.LoadOrStore(id, fc)
		if !loaded {
			break
		}
		v.(io.Closer).Close()
	}
	go cp.loopRecv(id, fc)
}

//loopCheck loop check conn done
func (cp *connpool) loopCheck() {
L:
	for {
		select {
		case <-cp.ctx.Done():
			break L
		case id := <-cp.rm:
			if v, loaded := cp.m.LoadAndDelete(id); loaded {
				v.(*fatconn).Cancel()
			}
		}
	}
	cp.m.Range(func(_, value interface{}) bool {
		cp.lookErr(value.(Conn).Close())
		return true
	})
	///this connpool will invalid after here
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
func (cp *connpool) loopRecv(id string, conn *fatconn) {
	defer conn.Close()
	for {
		d, err := conn.Recv()
		if err = cp.lookErr(err); err != nil {
			return
		}
		go cp.h.Process(id, d)
	}
}

//Send send given message
//does not block
func (cp *connpool) Send(v interface{}, tis ...string) (err error) {
	d, err := MakeDoc(v)
	if err != nil {
		return
	}
	if len(tis) == 0 {
		//send for all conn
		cp.m.Range(func(_, value interface{}) bool {
			cp.lookErr(value.(*fatconn).Send(d))
			return true
		})
		return
	}
	cp.m.Range(func(key, value interface{}) bool {
		for i, id := range tis {
			if id == key {
				copy(tis[i:], tis[i+1:])
				tis = tis[:len(tis)-1]
				if tmperr := cp.lookErr(value.(*fatconn).Send(d)); tmperr != nil {
					err = tmperr
				}
				break
			}
		}
		return len(tis) > 0
	})
	if len(tis) > 0 {
		err = fmt.Errorf("cannot found conn for (%v)", tis)
	}
	return
}

//fatconn package for basic Conn
//send in order
//close with status check
type fatconn struct {
	Cancel context.CancelFunc
	conn   Conn
	t      run.Tasker
	io.Closer
}

func newFatconn(ctx context.Context, id string, conn Conn, rm chan string) *fatconn {
	ctx, cancel := context.WithCancel(ctx)
	return &fatconn{
		Cancel: cancel,
		conn:   conn,
		t:      run.NewTasker(ctx, time.Second),
		Closer: run.NewOnceCloser(func() error {
			//when ctx closed, means connloop closed
			if ctx.Err() == nil {
				//notice connpool remove this conn
				rm <- id
				<-ctx.Done()
			}
			return conn.Close()
		}),
	}
}

//Write bytes to remote
func (fc *fatconn) Send(d *Doc) error {
	return fc.t.Add(run.Task{
		H: fc.handle,
		V: d,
	})
}

func (fc *fatconn) Recv() (*Doc, error) {
	return fc.conn.Recv()
}

//handle for Task
func (fc *fatconn) handle(v interface{}) (err error) {
	if werr := fc.conn.Send(v.(*Doc)); werr != nil {
		if werr != io.EOF {
			klog.Warningf("conn write error: %s", werr.Error())
		}
		//close conn and tell connpool to delete this conn
		fc.Close()
	}
	return
}
