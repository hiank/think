package run

import (
	"context"
	"io"
)

// type Task struct {
// 	//for handle value (V)
// 	H func(v any) error
// 	//value for handle (H)
// 	V any
// 	//for notice handle (H) error
// 	C chan error
// }

type Tasker interface {
	Add(Task) error
	io.Closer
}

type Task interface {
	Process() error
}

// type Tasker interface {
// 	Add(Task) error
// 	Stop()
// }

type Token interface {
	context.Context
	Cancel()
	Fork(ForkOptions) Token
}

type TokenSet interface {
	Get(uid uint64) Token
}

// type Handler interface {
// 	Handle(v any) error
// }
