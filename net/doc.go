package net

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Doc struct {
	b []byte
	v *anypb.Any
}

//MakeDoc make Doc
func MakeDoc(v interface{}) (d *Doc, err error) {
	d = new(Doc)
	switch v := v.(type) {
	case proto.Message:
		var ok bool
		if d.v, ok = v.(*anypb.Any); !ok {
			if d.v, err = anypb.New(v); err != nil {
				return nil, err
			}
		}
		d.b, err = proto.Marshal(d.v)
	case []byte:
		d.b, d.v = v, new(anypb.Any)
		err = proto.Unmarshal(d.b, d.v)
	default:
		err = ErrInvalidDocParam
	}
	if err != nil {
		d = nil
	}
	return
}

func (d *Doc) Bytes() []byte {
	return d.b
}

func (d *Doc) Any() *anypb.Any {
	return d.v
}

func (d *Doc) TypeName() string {
	return string(d.v.ProtoReflect().Descriptor().Name())
}

// type anyDoc struct {
// 	s    *anypb.Any
// 	d    []byte
// 	err  error
// 	once sync.Once
// }

// func (ad *anyDoc) Any() (*anypb.Any, error) {
// 	return ad.s, nil
// }

// func (ad *anyDoc) Bytes() ([]byte, error) {
// 	ad.once.Do(func() {
// 		ad.d, ad.err = proto.Marshal(ad.s)
// 	})
// 	return ad.d, ad.err
// }

// type bytesDoc struct {
// 	s    []byte
// 	d    *anypb.Any
// 	err  error
// 	once sync.Once
// }

// func (bd *bytesDoc) Bytes() ([]byte, error) {
// 	return bd.s, nil
// }

// func (bd *bytesDoc) Any() (*anypb.Any, error) {
// 	bd.once.Do(func() {
// 		amsg := new(anypb.Any)
// 		if bd.err = proto.Unmarshal(bd.s, amsg); bd.err == nil {
// 			bd.d = amsg
// 		}
// 	})
// 	return bd.d, bd.err
// }
