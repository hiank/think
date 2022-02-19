package file

import (
	"github.com/hiank/think/doc"
)

type Decoder interface {
	LoadFile(paths ...string) error
	LoadBytes(form Form, vals ...[]byte) error
	Decode(outVals ...interface{}) error
	Clear()
}

type Form int

const (
	FormInvalid Form = iota
	FormYaml
	FormJson
	FormRows
	FormGob
	FormPB
)

//Fit simple Buffer
//non-async safe
func Fit(form Form) Decoder {
	var b *doc.B
	switch form {
	case FormYaml:
		b = doc.Y.MakeB(nil)
	case FormJson:
		b = doc.J.MakeB(nil)
	case FormGob:
		b = doc.G.MakeB(nil)
	case FormPB:
		b = doc.P.MakeB(nil)
	case FormRows:
		b = doc.NewMaker(&doc.RowsCoder{RC: excelRowsReader(1)}).MakeB(nil)
	default: ///not support form
		return nil
	}
	return &fit{form: form, b: b}
}

func Fat() Decoder {
	return &fat{m: make(map[string]Decoder)} //new(fat)
}
