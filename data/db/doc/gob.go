package doc

import (
	"bytes"
	"encoding/gob"
)

type Gob []byte

func newGob(v []byte) Doc {
	gb := new(Gob)
	if v != nil {
		*gb = v
	}
	return gb
}

func (gb *Gob) Decode(v interface{}) error {
	return gob.NewDecoder(bytes.NewReader(*gb)).Decode(v)
}

func (gb *Gob) Encode(v interface{}) (err error) {
	b := new(bytes.Buffer)
	if err = gob.NewEncoder(b).Encode(v); err == nil {
		*gb = b.Bytes()
	}
	return
}

func (gb *Gob) Val() string {
	return string(*gb)
}
