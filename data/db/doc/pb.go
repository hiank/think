package doc

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

type PB []byte

func newPB(v []byte) Doc {
	pb := new(PB)
	if v != nil {
		*pb = v
	}
	return pb
}

func (pb *PB) Decode(v interface{}) (err error) {
	if msg, ok := v.(proto.Message); ok {
		err = proto.Unmarshal(*pb, msg)
	} else {
		err = errors.New("value to decode should be proto.Message")
	}
	return
}

func (pb *PB) Encode(v interface{}) (err error) {
	if msg, ok := v.(proto.Message); ok {
		var buf []byte
		if buf, err = proto.Marshal(msg); err == nil {
			*pb = PB(buf)
		}
	} else {
		err = errors.New("value to decode should be proto.Message")
	}
	return
}

func (pb *PB) Val() string {
	return string(*pb)
}
