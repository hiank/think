package net

// type Conn interface {
// 	Read() ([]byte, error)
// 	Write([]byte) error
// 	Close() error
// }

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

// type Connpool interface {
// 	//send message (use stored conn)
// 	Send(v interface{}, tis ...string) error
// 	//store new conn
// 	//when id existed, close old conn before
// 	AddConn(id string, c Conn)
// 	//add handler for handle message recved
// 	AddHandler(k interface{}, h Handler)
// 	//close all stored conn
// 	Close() error
// }

type Dialer interface {
	Dial(addr string) (IAC, error)
}

// type Client interface {
// 	Send(interface{}) error
// 	Close() error
// }

type Client interface {
	Send(d *Doc, ti string) error
	AddHandler()
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
	AddHandler(k interface{}, h Handler)
	//
	Close() error
}

//Handler handle message
type Handler interface {
	Handle(*Doc)
}

// //KeyDecoder decode value to key (string)
// type KeyDecoder interface {
// 	Decode(interface{}) string
// }

// type BytesCoder interface {
// 	Decode([]byte) (interface{}, error)
// 	Encode(interface{}) ([]byte, error)
// }

// //FuncKeyDecoder convert func to KeyDecoder interface
// type FuncKeyDecoder func(interface{}) string

// func (fkd FuncKeyDecoder) Decode(k interface{}) string {
// 	return fkd(k)
// }

// type CarrierHandler interface {
// 	Handle(*pb.Carrier)
// }

// type MessageHandler interface {
// 	Handle(id uint64, msg proto.Message)
// }

// //CarrierConverter convert Carrier to string key (use in HandleMux)
// //HandleMux use the converter to known which CarrierHandler registered use to Handle carrier message
// type CarrierConverter interface {
// 	GetKey(*pb.Carrier) (key string, ok bool)
// }

// var (
// 	//NewServer new a Server.
// 	//CarrierHandler use to handle message received from client (by Conn)
// 	NewServer func(Listener, CarrierHandler) Server = newServer

// 	NewClient func(Dialer) Client = newClient
// )

// type BytesCoder intera
