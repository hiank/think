package run

import (
	"io"
)

type Contextkey string

//Tasker sequential task executor
type Tasker interface {
	Add(Task) error
	io.Closer
	internalOnly()
}

type Task interface {
	Process() error
}

type Hooker[T any] interface {
	Hook(v T)
}

// type Token interface {
// 	context.Context
// 	Fork(...TokenOption) Token
// 	io.Closer
// 	// internal
// }

// type TokenSet interface {
// 	Get(uid uint64) Token
// }

// type internal interface {
// 	internalOnly()
// }

// type internalLimit struct{}

// func (internalLimit) internalOnly() {}
