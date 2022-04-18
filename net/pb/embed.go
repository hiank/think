package pb

import (
	"github.com/hiank/think/doc"
	"github.com/hiank/think/run"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	ErrNotProtobufValue = run.Err("pb: param for make M must be proto.Message|[]byte")
)

func toAny(m proto.Message) (out *anypb.Any, err error) {
	var ok bool
	if out, ok = m.(*anypb.Any); !ok {
		out, err = anypb.New(m)
	}
	return
}

//MakeM make M
func MakeM(v any) (m M, err error) {
	switch v := v.(type) {
	case proto.Message:
		if m.a, err = toAny(v); err == nil {
			m.b = doc.P.MakeB(nil)
			err = m.b.Encode(m.a)
		}
	case []byte:
		m.b, m.a = doc.P.MakeB(v), new(anypb.Any)
		err = m.b.Decode(m.a)
	default:
		err = ErrNotProtobufValue
	}
	return
}

type M struct {
	b *doc.B
	a *anypb.Any
}

func (m M) Bytes() []byte {
	return m.b.D
}

func (m M) Any() *anypb.Any {
	return m.a
}

func (m M) TypeName() (tn string) {
	if fn := m.a.MessageName(); fn.IsValid() {
		tn = string(fn.Name())
	}
	return
}
