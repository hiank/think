package doc

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"

	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

///********************Yaml under*******************///
type Yaml []byte

func (ym *Yaml) Encode(v interface{}) (err error) {
	buf, ok := v.([]byte)
	if !ok {
		if buf, err = yaml.Marshal(v); err != nil {
			return
		}
	}
	*ym = buf
	return
}

func (ym *Yaml) Decode(out interface{}) error {
	return yaml.Unmarshal(*ym, out)
}

func (ym *Yaml) Val() []byte {
	return *ym
}

///********************Json under*******************///
type Json []byte

func (js *Json) Encode(v interface{}) (err error) {
	buf, ok := v.([]byte)
	if !ok {
		if buf, err = json.Marshal(v); err != nil {
			return
		}
	}
	*js = buf
	return
}

func (js *Json) Decode(out interface{}) (err error) {
	return json.Unmarshal(*js, out)
}

func (js *Json) Val() []byte {
	return *js
}

///********************Protobuf under*******************///
type PB []byte

func (pb *PB) Encode(v interface{}) (err error) {
	buf, ok := v.([]byte)
	if !ok {
		msg, ok := v.(proto.Message)
		if !ok {
			return errors.New("value to decode should be proto.Message")
		}
		if buf, err = proto.Marshal(msg); err != nil {
			return err
		}
	}
	*pb = buf
	return
}

func (pb *PB) Decode(out interface{}) (err error) {
	if msg, ok := out.(proto.Message); ok {
		err = proto.Unmarshal(*pb, msg)
	} else {
		err = errors.New("value to decode should be proto.Message")
	}
	return
}

func (pb *PB) Val() []byte {
	return *pb
}

///********************Gob under*******************///
type Gob []byte

func (gb *Gob) Encode(v interface{}) (err error) {
	buf, ok := v.([]byte)
	if !ok {
		b := new(bytes.Buffer)
		if err = gob.NewEncoder(b).Encode(v); err != nil {
			return
		}
		buf = b.Bytes()
	}
	*gb = buf
	return
}

func (gb *Gob) Decode(v interface{}) error {
	return gob.NewDecoder(bytes.NewReader(*gb)).Decode(v)
}

func (gb *Gob) Val() []byte {
	return *gb
}
