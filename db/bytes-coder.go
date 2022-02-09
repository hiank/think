package db

import (
	"fmt"

	"github.com/hiank/think/doc"
)

//PB protobuf type value
type PB struct {
	V interface{}
}

//GOB gob type value
type GOB struct {
	V interface{}
}

//JSON json type value
type JSON struct {
	V interface{}
}

type BytesCoder byte

//Encode encode value to bytes
func (bc BytesCoder) Encode(v interface{}) (out []byte, err error) {
	var d doc.Doc
	if pb, ok := v.(PB); ok {
		d, v = new(doc.PB), pb.V
	} else if gb, ok := v.(GOB); ok {
		d, v = new(doc.Gob), gb.V
	} else if js, ok := v.(JSON); ok {
		d, v = new(doc.Json), js.V
	} else {
		return nil, fmt.Errorf("invalid param type: support PB GOB JSON now")
	}
	if err = d.Encode(v); err == nil {
		out = d.Val()
	}
	return
}

//Decode decode bytes to out value
func (bc BytesCoder) Decode(v []byte, out interface{}) error {
	if pb, ok := out.(PB); ok {
		d := doc.PB(v)
		return d.Decode(pb.V)
	}
	if gb, ok := out.(GOB); ok {
		d := doc.Gob(v)
		return d.Decode(gb.V)
	}
	if js, ok := out.(JSON); ok {
		d := doc.Json(v)
		return d.Decode(js.V)
	}
	return fmt.Errorf("invalid param type: support PB GOB JSON now")
}
