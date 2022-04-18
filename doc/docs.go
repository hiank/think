package doc

import (
	"fmt"
)

type maker struct {
	coder Coder
}

func NewMaker(coder Coder) Maker {
	return &maker{coder: coder}
}

func (m *maker) MakeT(v any) T {
	return T{
		V:          v,
		embedCoder: embedCoder{m.coder},
	}
}

func (m *maker) MakeB(d []byte) *B {
	return &B{
		D:          d,
		embedCoder: embedCoder{m.coder},
	}
}

type embedCoder struct {
	c Coder
}

func (ec embedCoder) Decode(data []byte, out any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	return ec.c.Decode(data, out)
}

func (ec embedCoder) Encode(v any) (out []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	return ec.c.Encode(v)
}

// type

type T struct {
	V any
	// embed Coder
	embedCoder
}

//Decode V form data
func (t T) Decode(data []byte) error {
	return t.embedCoder.Decode(data, t.V)
}

//Encode V to data
func (t T) Encode() ([]byte, error) {
	return t.embedCoder.Encode(t.V)
}

type B struct {
	D []byte
	embedCoder
}

//Decode out from D
func (b B) Decode(out any) (err error) {
	return b.embedCoder.Decode(b.D, out)
}

//Encode v to D
func (b *B) Encode(v any) (err error) {
	if d, ok := v.([]byte); ok {
		b.D = d
	} else {
		b.D, err = b.embedCoder.Encode(v)
	}
	return
}
