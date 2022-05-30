package one

import "github.com/hiank/think/net/box"

type Tokenset interface {
	// //Build build a new Token (root Token) and return it
	// //when root Token existed, return error
	// Build(uid string) (box.Token, error)

	//Derive get a derived Token
	//when (non root Token)/(root Token canceled), return error
	Derive(uid string) box.Token

	//Kill cancel root Token
	//when non root Token, return error
	Kill(uid string) error
}
