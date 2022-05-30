package net

import (
	"context"
	"io"

	"github.com/hiank/think/net/box"
	"github.com/hiank/think/net/one"
	"github.com/hiank/think/run"
	"google.golang.org/protobuf/proto"
)

const (
	DefaultHandler string = ""
)

type server struct {
	ctx      context.Context
	listener Listener
	cp       *connset
	io.Closer
}

func NewServer(listener Listener, h Handler) Server {
	srv := &server{
		listener: listener,
		cp:       newConnset(h),
	}
	srv.ctx, srv.Closer = run.StartHealthyMonitoring(one.TODO(), run.CloserToDoneHook(listener), srv.cp.close)
	return srv
}

//ListenAndServe block to accept new conn until the listener closed or server closed
func (srv *server) ListenAndServe() error {
	defer srv.Close()
	for {
		tc, err := srv.listener.Accept()
		if err == nil {
			if err = srv.ctx.Err(); err == nil {
				srv.cp.loadOrStore(srv.ctx, tc.Token.Value(box.ContextkeyTokenUid).(string), func(context.Context) (TokenConn, error) {
					return tc, nil
				})
				continue
			}
		}
		return err
	}
}

func (srv *server) Send(pm proto.Message, tis ...string) (err error) {
	select {
	case <-srv.ctx.Done():
		err = srv.ctx.Err()
	default:
		m := box.New(box.WithMessageValue(pm))
		switch len(tis) {
		case 0:
			err = srv.cp.broadcast(m)
		default:
			err = srv.cp.multiSend(m, tis...)
		}
	}
	return
}
