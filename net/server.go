package net

import (
	"context"
	"reflect"
	"sync"

	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/run"
	"k8s.io/klog/v2"
)

const (
	DefaultHandler string = ""
)

type server struct {
	ctx      context.Context
	cancel   context.CancelFunc
	listener Listener
	*connpool
}

func NewServer(listener Listener, h Handler) Server {
	ctx, cancel := context.WithCancel(run.TODO())
	return &server{
		listener: listener,
		ctx:      ctx,
		cancel:   cancel,
		connpool: newConnpool(ctx, h),
	}
}

//ListenAndServe block to accept new conn until the listener closed or server closed
func (srv *server) ListenAndServe() (err error) {
	defer srv.Close()
	for {
		iac, err := srv.listener.Accept()
		if err == nil {
			if err = srv.ctx.Err(); err == nil {
				srv.AddConn(iac.ID, iac.Conn)
				continue
			}
		}
		return err
	}
}

//Close close the server
//will close all conns then clear the conns's map
//the method could be called multiple
func (srv *server) Close() (err error) {
	if err = srv.ctx.Err(); err == nil {
		srv.cancel() //will clean connpool by this call
		err = srv.listener.Close()
	}
	return
}

type RouteMux struct {
	m sync.Map
}

func (rm *RouteMux) Handle(k interface{}, h Handler) {
	sk, ok := k.(string)
	if !ok {
		rv := reflect.ValueOf(k)
		for rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		sk = rv.Type().Name()
	}
	rm.m.Store(sk, h)
}

func (rm *RouteMux) Route(id string, m pb.M) {
	mv, loaded := rm.m.Load(m.TypeName())
	if !loaded {
		if mv, loaded = rm.m.Load(DefaultHandler); !loaded {
			klog.Warning("cannot find handler for handle message recv by conn: ", m.TypeName())
			return
		}
	}
	mv.(Handler).Route(id, m)
}
