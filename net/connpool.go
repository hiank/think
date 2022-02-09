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
	ctx  context.Context
	h    Handler    //Handler for revc message
	m    sync.Map   //conn map
	errc chan error //for recv write error conn's identity
}

func newConnpool(ctx context.Context, h Handler) (cp *connpool) {
	cp = &connpool{
		ctx:  ctx,
		h:    h,
		errc: make(chan error),
	}
	go cp.loopCheck()
	return
}

//AddConn add new conn to pool
//when existed same id, delete it and close it's conn before
//wait unit add success
func (cp *connpool) AddConn(id string, c Conn) {
	v, loaded := cp.m.LoadOrStore(id, newFatconn(cp.ctx, id, c, cp.errc))
	if !loaded {
		///conn stored sucess
		///loop recv from conn here
		go cp.loopRecv(id, v.(*fatconn))
		return
	}
	//wait until stored conn deleted
	<-v.(*fatconn).Done()
	cp.AddConn(id, c)
}

//loopCheck loop check conn done
func (cp *connpool) loopCheck() {
L:
	for {
		select {
		case <-cp.ctx.Done():
			break L
		case err := <-cp.errc:
			if v, loaded := cp.m.LoadAndDelete(err.Error()); loaded {
				conn := v.(*fatconn)
				conn.Cancel()
				conn.Close()
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
func (cp *connpool) loopRecv(identity string, conn *fatconn) {
	defer conn.Done()
	for {
		d, err := conn.Recv()
		if err = cp.lookErr(err); err != nil {
			return
		}
		// if v, err := cp.coder.Encode(d); err == nil {
		go cp.h.Handle(d)
		// } else {
		// 	//print error. err should be encode error(should not be io.EOF/context closed)
		// 	cp.lookErr(err)
		// }
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

//Close close connpool
//Deprecated: close ctx instead of call Close. the resources will be cleaned up automaic after ctx closed
func (cp *connpool) Close() error {
	return fmt.Errorf("please close ctx instead of Close. after ctx closed, the conn will be clear automatic")
}

//fatconn package for basic Conn
//send in order
//close with status check
type fatconn struct {
	ctx    context.Context
	Cancel context.CancelFunc
	Conn
	id   string
	t    run.Tasker
	errc chan error
	once sync.Once
}

func newFatconn(ctx context.Context, id string, conn Conn, errc chan error) *fatconn {
	ctx, cancel := context.WithCancel(ctx)
	return &fatconn{
		ctx:    ctx,
		Cancel: cancel,
		id:     id,
		Conn:   conn,
		errc:   errc,
		t:      run.NewTasker(ctx, time.Second),
	}
}

//Write bytes to remote
func (fc *fatconn) Send(d *Doc) error {
	return fc.t.Add(run.Task{
		H: fc.handle,
		V: d,
	})
}

//Done close and wait until conn closed
func (fc *fatconn) Done() <-chan struct{} {
	fc.once.Do(func() {
		if fc.ctx.Err() == nil {
			//ctx not closed, notice connpool to close and delete this conn
			fc.errc <- fmt.Errorf("%s", fc.id)
		}
	})
	return fc.ctx.Done()
}

//handle for Task
func (fc *fatconn) handle(v interface{}) (err error) {
	if werr := fc.Conn.Send(v.(*Doc)); werr != nil {
		if werr != io.EOF {
			klog.Warningf("conn for id (%s) write error: %s", fc.id, werr.Error())
		}
		//tell connpool the conn for (id) write error
		//the id will send by fc.errc (recv in connpool's loopCheck)
		fc.Done()
	}
	return
}
