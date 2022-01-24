package ws

import (
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hiank/think/net"
	"github.com/hiank/think/oauth"
	"k8s.io/klog/v2"
)

type listener struct {
	srv      *http.Server
	pp       chan net.Conn
	upgrader *websocket.Upgrader
	auther oauth.Auther
}

func NewListener(auther oauth.Auther, addr string) net.Listener {
	l := &listener{
		pp:       make(chan net.Conn),
		upgrader: &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
		auther:   auther,
	}
	httpHandler := http.NewServeMux()
	httpHandler.Handle("/ws", l)
	l.srv = &http.Server{Addr: addr, Handler: httpHandler}

	go func(srv *http.Server, pp chan<- net.Conn) {
		klog.Warning("websocket server stopped: ", srv.ListenAndServe())
		close(pp)
	}(l.srv, l.pp)
	return l
}

func (l *listener) Accept() (c net.Conn, err error) {
	c, ok := <-l.pp
	if !ok {
		err = io.EOF
	}
	return c, err
}

func (l *listener) Close() error {
	return l.srv.Close()
}

func (l *listener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	if token == "" {
		http.Error(w, "non token value in request header", http.StatusNonAuthoritativeInfo)
		return
	}
	uid, err := l.auther.Auth(token)
	if err != nil {
		http.Error(w, "token invalid: "+err.Error(), http.StatusUnauthorized)
		return
	}

	wc, err := l.upgrader.Upgrade(w, r, nil)
	if err != nil {
		klog.Warning("ws: Upgrade error: ", err)
		return
	}
	l.pp <- &conn{wc: wc, uid: uid}
}
