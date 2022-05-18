package ws

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/hiank/think/net"
	"github.com/hiank/think/net/adapter"
	"github.com/hiank/think/oauth"
	"github.com/hiank/think/run"

	"k8s.io/klog/v2"
)

type listener struct {
	upgrader *websocket.Upgrader
	auther   oauth.Auther
	io.Closer
	adapter.ChanAccepter
}

func (lis *listener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	if token == "" {
		http.Error(w, "non token value in request header", http.StatusNonAuthoritativeInfo)
		return
	}
	uid, err := lis.auther.Auth(token)
	if err != nil {
		http.Error(w, "token invalid: "+err.Error(), http.StatusUnauthorized)
		return
	}

	wc, err := lis.upgrader.Upgrade(w, r, nil)
	if err != nil {
		klog.Warning("ws: Upgrade error: ", err)
		return
	}
	lis.ChanAccepter <- net.IdentityConn{ID: strconv.FormatUint(uid, 10), Conn: &conn{wc: wc}}
}

func (lis *listener) contextHealthy(ctx context.Context, lisCloser io.Closer) {
	<-ctx.Done()
	close(lis.ChanAccepter)
	lisCloser.Close()
}

func NewListener(ctx context.Context, opt ListenOption) net.Listener {
	ctx, cancel := context.WithCancel(ctx)
	opt = withDefaultListenOption(opt)
	lis := &listener{
		ChanAccepter: make(adapter.ChanAccepter),
		upgrader:     &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
		auther:       opt.Auther,
		Closer: run.NewOnceCloser(func() error {
			cancel()
			return nil
		}),
	}
	httpHandler := http.NewServeMux()
	httpHandler.Handle("/ws", lis)
	srv := &http.Server{Addr: opt.Addr, Handler: httpHandler}

	go func() {
		klog.Warning("websocket server stopped: ", srv.ListenAndServe())
	}()
	go lis.contextHealthy(ctx, srv)
	return lis
}
