package auth

import (
	"context"
	"io"
)

type Tokenset interface {
	// //Build build a new Token (root Token) and return it
	// //when root Token existed, return error
	// Build(uid string) (box.Token, error)

	//Derive get a derived Token
	//when (non root Token)/(root Token canceled), return error
	Derive(key string, topts ...TokenOption) Token

	//Kill cancel root Token
	//when non root Token, return error
	Kill(key string) error

	io.Closer

	//could only defined internal
	internalOnly()
}

type Token interface {
	context.Context
	Fork(...TokenOption) Token
	ToString() string //token string value
	io.Closer

	//could only defined internal
	internalOnly()
}

type internal struct{}

func (internal) internalOnly() {}
