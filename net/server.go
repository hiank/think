package net

import (
	"context"
	"sync"

	"github.com/hiank/think/net/box"
	"github.com/hiank/think/run"
	"google.golang.org/protobuf/proto"
	"k8s.io/klog/v2"
)

const (
	DefaultHandler string = ""
)

type server struct {
	ctx      context.Context
	cancel   context.CancelFunc
	listener Listener
	cp       *connpool
}

func NewServer(listener Listener, h Handler) Server {
	ctx, cancel := context.WithCancel(run.TODO())
	return &server{
		listener: listener,
		ctx:      ctx,
		cancel:   cancel,
		cp:       newConnpool(ctx, h),
	}
}

//ListenAndServe block to accept new conn until the listener closed or server closed
func (srv *server) ListenAndServe() (err error) {
	defer srv.Close()
	for {
		iac, err := srv.listener.Accept()
		if err == nil {
			if err = srv.ctx.Err(); err == nil {
				srv.cp.add(iac.ID, iac.Conn)
				continue
			}
		}
		return err
	}
}

func (srv *server) Send(pm proto.Message, tis ...string) (err error) {
	m, err := box.New(pm)
	if err == nil {
		switch len(tis) {
		case 0:
			err = srv.cp.broadcast(m)
		default:
			err = srv.cp.multiSend(m, tis...)
		}
	}
	return
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
	}
	rm.m.Store(sk, h)
}

func (rm *RouteMux) Route(id string, m *box.Message) {
	k := string(m.GetAny().MessageName().Name())
	mv, loaded := rm.m.Load(k)
	if !loaded {
		if mv, loaded = rm.m.Load(DefaultHandler); !loaded {
			klog.Warning("cannot find handler for handle message recv by conn: ", k)
			return
		}
	}
	mv.(Handler).Route(id, m)
}
