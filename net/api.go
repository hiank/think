package net

type Conn interface {
	Send(*Doc) error
	Recv() (*Doc, error)
	Close() error
}

//IAC identity and conn
type IAC struct {
	ID string
	Conn
}

//Dialer dial to server
type Dialer interface {
	Dial(addr string) (IAC, error)
}

type Client interface {
	Send(d *Doc, ti string) error
	Handle(k interface{}, h Handler)
}

type Listener interface {
	Accept() (IAC, error)
	Close() error
}

type Server interface {
	//start work
	ListenAndServe() error
	//Send message to client (by conn)
	//ti is target identity. when len(ti) == 0
	//means send for all conn
	Send(v interface{}, tis ...string) error
	//AddHandler add handler for revc message
	Handle(k interface{}, h Handler)
	//
	Close() error
}

//Handler handle message
type Handler interface {
	Process(string, *Doc)
}
