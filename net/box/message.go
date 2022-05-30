package box

import (
	"github.com/hiank/think/doc"
	"github.com/hiank/think/run"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/klog/v2"
)

type message struct {
	a anypb.Any
	b []byte
	doc.PBCoder
}

func (msg *message) GetAny() *anypb.Any {
	return &msg.a
}

func (msg *message) GetBytes() []byte {
	return msg.b
}

func (msg *message) internalOnly() {}

//Unmarhsal data to out. out must be *message bacause limit by internalOnly
func Unmarshal[T proto.Message](data []byte, out Message) (err error) {
	tv, err := doc.MakeT[T]()
	if err == nil {
		if err = out.Decode(data, tv); err == nil {
			err = unmarshalProtoMessage(tv, out.(*message))
		}
	}
	return
}

//UnmarshalNew unmarshal bytes to new *Message.
func UnmarshalNew[T proto.Message](data []byte) (out Message, err error) {
	msg := &message{}
	if err = Unmarshal[T](data, msg); err == nil {
		out = msg
	}
	return
}

func unmarshalProtoMessage(pm proto.Message, out *message) (err error) {
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

type MessageOption run.Option[*message]

func WithMessageValue(pm proto.Message) MessageOption {
	return run.FuncOption[*message](func(msg *message) {
		if err := unmarshalProtoMessage(pm, msg); err != nil {
			klog.Warning("box: failed to new message", err)
		}
	})
}

//New message
//non MessageOption return empty message
func New(opts ...MessageOption) Message {
	msg := &message{}
	for _, opt := range opts {
		opt.Apply(msg)
	}
	return msg
}
