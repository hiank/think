package pb_test

import (
	"testing"

	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/set/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

func TestProtoMessageName(t *testing.T) {

	msg := &pb.Message{Key: "WS"}
	assert.Equal(t, "Message", string(msg.ProtoReflect().Descriptor().Name()))
}

func TestProtoMessageNameInAny(t *testing.T) {
	anyMsg, _ := anypb.New(&pb.Message{Key: "WS"})
	name := anyMsg.MessageName().Name()
	assert.Assert(t, name.IsValid())

	assert.Equal(t, string(name), "Message")

	t.Run("Empty", func(t *testing.T) {
		anyMsg = new(anypb.Any)
		assert.Assert(t, !anyMsg.MessageName().IsValid())
		name := anyMsg.MessageName().Name()
		assert.Assert(t, !name.IsValid())
	})
}

type testHandler struct {
	out chan<- proto.Message
}

func (th *testHandler) Handle(msg proto.Message) error {
	th.out <- msg
	return nil
}

func TestDefaultHandler(t *testing.T) {

	out := make(chan proto.Message, 1)
	t.Run("Register", func(t *testing.T) {
		err := pb.DefaultHandler.Register(new(pb.Message), &testHandler{out: out})
		assert.Assert(t, err == nil)

		err = pb.DefaultHandler.Register(new(pb.Message), &testHandler{out: out})
		assert.Equal(t, err.Error(), codes.ErrorExistedMessageHandler.Error())
	})

	t.Run("HandleMessage", func(t *testing.T) {
		pb.DefaultHandler.Handle(&pb.Message{Key: "normalMessage"})
		msg := <-out
		assert.Equal(t, msg.(*pb.Message).GetKey(), "normalMessage")
	})

	t.Run("HandleAny", func(t *testing.T) {
		anyMsg, _ := anypb.New(&pb.Message{Key: "anyMessage"})
		pb.DefaultHandler.Handle(anyMsg)
		msg := <-out
		assert.Equal(t, msg.(*pb.Message).GetKey(), "anyMessage")
	})
}
