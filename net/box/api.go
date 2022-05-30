package box

import (
	"context"
	"io"

	"github.com/hiank/think/doc"
	"google.golang.org/protobuf/types/known/anypb"
)

type Token interface {
	context.Context
	Fork(...TokenOption) Token
	io.Closer
	internalOnly()
}

type Message interface {
	GetAny() *anypb.Any
	GetBytes() []byte
	// doc.PBCoder
	doc.Coder
	internalOnly()
}

type TT[T any] struct {
	Token Token
	T     T
}
