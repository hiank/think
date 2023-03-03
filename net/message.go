package net

import (
	"sync"

	"github.com/hiank/think/auth"
	"github.com/hiank/think/doc"
	"github.com/hiank/think/run"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	ErrNonMessageValue = run.Err("net: at least one of WithMessageBytes/WithMessageValue needed for NewMessage")
)

type Message struct {
	tk auth.Token
	a  *anypb.Any
	doc.Proto

	onceSet sync.Once //once set value
}

func (msg *Message) Token() auth.Token {
	return msg.tk
}

func (msg *Message) Any() *anypb.Any {
	return msg.a
}

type MessageOption run.Option[*Message]

// WithMessageToken set given Token to message
func WithMessageToken(tk auth.Token) MessageOption {
	return run.FuncOption[*Message](func(msg *Message) {
		msg.tk = tk
	})
}

// WithMessageValue value for encode
// when v is not a *anypb.Any, packing it to *anypb.Any first
func WithMessageValue(v proto.Message) MessageOption {
	return run.FuncOption[*Message](func(msg *Message) {
		msg.onceSet.Do(func() {
			var err error
			a, ok := v.(*anypb.Any)
			if !ok {
				if a, err = anypb.New(v); err != nil {
					panic(err)
				}
			}
			if err = msg.Encode(a); err != nil {
				panic(err)
			}
			msg.a = a
		})
	})
}

// WithMessageBytes bytes for decode to *anypb.Any
// it would panic when given bts was not *anypb.Any buff
func WithMessageBytes(bts []byte) MessageOption {
	return run.FuncOption[*Message](func(msg *Message) {
		msg.onceSet.Do(func() {
			msg.Proto, msg.a = bts, new(anypb.Any)
			if err := msg.Decode(msg.a); err != nil {
				panic(err)
			}
		})
	})
}

// NewMessage new a not empty message
func NewMessage(opts ...MessageOption) (m *Message) {
	m = &Message{}
	for _, opt := range opts {
		opt.Apply(m)
	}
	m.onceSet.Do(func() {
		panic(ErrNonMessageValue)
	})
	return
}
