package run

import (
	"io"
)

type Contextkey string

type Tasker interface {
	Add(Task) error
	io.Closer
	internalOnly()
}

type Task interface {
	Process() error
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
