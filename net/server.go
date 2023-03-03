package net

import (
	"context"
	"io"

	"github.com/hiank/think/run"
	"google.golang.org/protobuf/proto"
)

const (
	DefaultHandler string = "__default_handler_routemux__"
)

type Server struct {
	ctx      context.Context
	listener Listener
	cp       *connset
	io.Closer
}

func NewServer(ctx context.Context, lis Listener, h Handler) *Server {
	srv := &Server{
		listener: lis,
		cp:       newConnset(h),
	}
	srv.ctx, srv.Closer = run.StartHealthyMonitoring(ctx, run.CloserToDoneHook(lis), srv.cp.close)
	return srv
}

// ListenAndServe block to accept new conn until the listener closed or server closed
func (srv *Server) ListenAndServe() error {
	defer srv.Close()
	for {
		tc, err := srv.listener.Accept()
		if err == nil {
			if err = srv.ctx.Err(); err == nil {
				srv.cp.loadOrStore(srv.ctx, tc.Token().ToString(), func(context.Context) (Conn, error) {
					return tc, nil
				})
				continue
			}
		}
		return err
	}
}

func (srv *Server) Send(pm proto.Message, tis ...string) (err error) {
	if err = srv.ctx.Err(); err == nil {
		m := NewMessage(WithMessageValue(pm))
		switch len(tis) {
		case 0:
			err = srv.cp.broadcast(m)
		default:
			err = srv.cp.multiSend(m, tis...)
		}
	}
	return
}
