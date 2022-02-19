package run

import "context"

type Task struct {
	//for handle value (V)
	H func(v interface{}) error
	//value for handle (H)
	V interface{}
	//for notice handle (H) error
	C chan error
}

type Tasker interface {
	Add(Task) error
	Stop()
}

type Token interface {
	context.Context
	Cancel()
	Fork(ForkOptions) Token
}

type TokenSet interface {
	Get(uid uint64) Token
}

// type Handler interface {
// 	Handle(v interface{}) error
// }
