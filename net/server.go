package net

import (
	"context"
)

// func defaultServeOptions() serveOptions {
// 	return serveOptions{
// 		// handlerKeyDecoder: coder.TypeKey(0),
// 		// bytesCoder:        coder.AnyBytes(0),
// 	}
// }

type server struct {
	ctx      context.Context
	cancel   context.CancelFunc
	listener Listener
	*fathandler
	*connpool
}

func NewServer(listener Listener) Server {
	// dopts := defaultServeOptions()
	// for _, opt := range opts {
	// 	opt.apply(&dopts)
	// }
	ctx, cancel := context.WithCancel(context.Background())
	// h := &fathandler{kd: dopts.handlerKeyDecoder}
	h := new(fathandler)
	return &server{
		listener:   listener,
		ctx:        ctx,
		cancel:     cancel,
		fathandler: h,
		connpool:   newConnpool(ctx, h),
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
