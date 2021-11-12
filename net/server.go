package net

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/hiank/think/net/pb"
	"google.golang.org/protobuf/proto"
	"k8s.io/klog/v2"
)

// var (
// 	//DefaultHandleKey use for handleMux
// 	DefaultHandlerKey string = "_default_key_"
// )

func defaultHandleOptions() handleOptions {
	return handleOptions{
		converter: FuncCarrierConverter(func(carrier *pb.Carrier) (string, bool) {
			return string(carrier.GetMessage().MessageName().Name()), true
		}),
	}
}

type FuncCarrierConverter func(*pb.Carrier) (string, bool)

func (fcc FuncCarrierConverter) GetKey(carrier *pb.Carrier) (string, bool) {
	return fcc(carrier)
}

type handleMux struct {
	m     sync.Map
	dopts handleOptions
}

func NewHandleMux(opts ...HandleOption) *handleMux {
	hm := &handleMux{
		dopts: defaultHandleOptions(),
	}
	for _, opt := range opts {
		opt.apply(&hm.dopts)
	}
	return hm
}

//Look register handler for key
func (hm *handleMux) Look(key string, handler IMessageHandler) {
	hm.m.Store(key, handler)
}

//LookObject register handler by proto.Message instance
func (hm *handleMux) LookObject(obj proto.Message, handler IMessageHandler) {
	hm.Look(string(obj.ProtoReflect().Descriptor().Name()), handler)
}

//Handle handle given carrier message
//the method will find suitable handler to handle the message
func (hm *handleMux) Handle(carrier *pb.Carrier) {
	key, ok := hm.dopts.converter.GetKey(carrier)
	if ok {
		if val, ok := hm.m.Load(key); ok {
			handler, _ := val.(IMessageHandler)
			if msg, err := carrier.GetMessage().UnmarshalNew(); err == nil {
				handler.Handle(carrier.GetIdentity(), msg)
			} else {
				klog.Warning(err) //NOTE: unmarshal error
			}
			return
		}
		if hm.dopts.defaultHandler != nil {
			hm.dopts.defaultHandler.Handle(carrier)
			return
		}
	}
	klog.Warningf("cannot find handler to handle message (%s)", key)
}

type server struct {
	ctx      context.Context
	cancel   context.CancelFunc
	listener IListener
	handler  ICarrierHandler
	m        sync.Map
}

func newServer(listener IListener, handler ICarrierHandler) IServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &server{
		listener: listener,
		ctx:      ctx,
		cancel:   cancel,
		handler:  handler,
	}
}

//ListenAndServe block to accept new conn until the listener closed or server closed
func (srv *server) ListenAndServe() (err error) {
	defer srv.Close()
	for {
		conn, err := srv.listener.Accept()
		if err = srv.lookErr(err); err != nil {
			return err
		}
		go srv.handleConn(conn)
	}
}

//Send send message (with identity) to client (use conn)
//carrier's identity is target's value
func (srv *server) Send(carrier *pb.Carrier) (err error) {
	if val, ok := srv.m.Load(carrier.GetIdentity()); ok {
		err = val.(IConn).Send(carrier.GetMessage())
	} else {
		err = fmt.Errorf("cannot found conn (identity:%x) in the server", carrier.GetIdentity())
	}
	return
}

//Close close the server
//will close all conns then clear the conns's map
//the method could be called multiple
func (srv *server) Close() (err error) {
	if err = srv.ctx.Err(); err == nil {
		srv.cancel()
		err = srv.listener.Close()
		srv.m.Range(func(key, value interface{}) bool {
			value.(IConn).Close()
			return true
		})
	}
	return
}

//handleConn loop recv message by conn and handle it until conn closed or server closed
func (srv *server) handleConn(conn IConn) {
	defer conn.Close()
	identity := conn.GetIdentity()
	if _, loaded := srv.m.LoadOrStore(identity, conn); loaded {
		//identity would loaded already
		klog.Warningf("conn (identity:%x) already loaded", identity)
		return
	}

	defer srv.m.Delete(identity)
	for {
		any, err := conn.Recv()
		if err = srv.lookErr(err); err != nil {
			if err != io.EOF {
				klog.Warningf("conn (identity:%x) closed abnormally: %x", identity, err)
			}
			return
		}
		go srv.handler.Handle(&pb.Carrier{Identity: identity, Message: any})
	}
}

//lookErr check wether the ctx would canceled
func (srv *server) lookErr(err error) error {
	if srv.ctx.Err() != nil {
		err = srv.ctx.Err()
	}
	return err
}
