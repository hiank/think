package box

import (
	"github.com/hiank/think/doc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Message struct {
	a anypb.Any
	b []byte
	doc.PBCoder
}

func (msg *Message) GetAny() *anypb.Any {
	return &msg.a
}

func (msg *Message) GetBytes() []byte {
	return msg.b
}

func Unmarshal[T proto.Message](data []byte, out *Message) (err error) {
	tv, err := doc.MakeT[T]()
	if err == nil {
		if err = out.Decode(data, tv); err == nil {
			err = unmarshalProtoMessage(tv, out)
		}
	}
	return
}

//UnmarshalNew unmarshal bytes to new *Message.
func UnmarshalNew[T proto.Message](data []byte) (out *Message, err error) {
	out = new(Message)
	err = Unmarshal[T](data, out)
	return
}

func unmarshalProtoMessage(pm proto.Message, out *Message) (err error) {
	switch amsg := pm.(type) {
	case *anypb.Any:
		out.a = *amsg
	default:
		err = out.GetAny().MarshalFrom(pm)
	}
	if err == nil {
		out.b, err = out.Encode(out.GetAny())
	}
	return
}

func New(pm proto.Message) (out *Message, err error) {
	out = &Message{}
	err = unmarshalProtoMessage(pm, out)
	return
}
