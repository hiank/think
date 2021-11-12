package ws

import (
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hiank/think/net"
	"k8s.io/klog/v2"
)

type listener struct {
	srv      *http.Server
	pp       chan net.IConn
	upgrader *websocket.Upgrader
	storage  IStorage
}

func NewListener(storage IStorage, addr string) net.IListener {
	l := &listener{
		pp:       make(chan net.IConn),
		upgrader: &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
		storage:  storage,
	}
	httpHandler := http.NewServeMux()
	httpHandler.Handle("/ws", l)
	l.srv = &http.Server{Addr: addr, Handler: httpHandler}

	go func(srv *http.Server, pp chan<- net.IConn) {
		klog.Warning("websocket server stopped: ", srv.ListenAndServe())
		close(pp)
	}(l.srv, l.pp)
	return l
}

func (l *listener) Accept() (c net.IConn, err error) {
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
	uid, ok := l.storage.GetUidByToken(token)
	if !ok {
		http.Error(w, "token invalid", http.StatusUnauthorized)
		return
	}

	wc, err := l.upgrader.Upgrade(w, r, nil)
	if err != nil {
		klog.Warning("ws: Upgrade error: ", err)
		return
	}
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		klog.Warning("ws: accept chan closed: ", r)
	// 	}
	// }()
	l.pp <- &conn{wc: wc, uid: uid}
}
