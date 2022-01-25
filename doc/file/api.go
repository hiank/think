package file

import (
	"github.com/hiank/think/doc"
)

type Decoder interface {
	LoadFile(paths ...string) error
	LoadBytes(form Form, vals ...[]byte) error
	Decode(outVals ...interface{}) error
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
	// fit := &fit{form: form}
	var d doc.Doc
	switch form {
	case FormYaml:
		d = doc.YamlMaker.Make(nil)
	case FormJson:
		d = doc.JsonMaker.Make(nil)
	case FormGob:
		d = doc.GobMaker.Make(nil)
	case FormPB:
		d = doc.PBMaker.Make(nil)
	case FormRows:
		d = doc.NewRows(excelRowsReader(0))
	default: ///not support form
		return nil
	}
	return &fit{form: form, doc: d}
}

func Fat() Decoder {
	return new(fat)
}
