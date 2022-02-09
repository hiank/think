package testdata_test

import (
	"testing"

	"github.com/hiank/think/net/testdata"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

func TestMessageName(t *testing.T) {

	t.Run("G_Example", func(t *testing.T) {
		gmsg := &testdata.G_Example{Value: "g"}
		assert.Equal(t, string(gmsg.ProtoReflect().Descriptor().FullName()), "G_Example", "")
		anyMsg, _ := anypb.New(gmsg)
		assert.Equal(t, anyMsg.MessageName(), gmsg.ProtoReflect().Descriptor().FullName(), "")

		msg, _ := anyMsg.UnmarshalNew()
		assert.Equal(t, msg.ProtoReflect().Descriptor().FullName(), anyMsg.MessageName(), "")
	})

	t.Run("Test1", func(t *testing.T) {
		tmsg := &testdata.Test1{Name: "t"}
		// t.Log(tmsg.ProtoReflect().Descriptor().Name())
		assert.Equal(t, string(tmsg.ProtoReflect().Descriptor().FullName()), "test1", "the name is the message name in .proto")
		anyMsg, _ := anypb.New(tmsg)
		assert.Equal(t, tmsg.ProtoReflect().Descriptor().FullName(), anyMsg.MessageName(), "name is same by get from any message and normal message")

		msg, _ := anyMsg.UnmarshalNew()
		assert.Equal(t, msg.ProtoReflect().Descriptor().FullName(), anyMsg.MessageName())
		// msg, _ := anyMsg.UnmarshalNew()
		// assert.Equal(t,
	})
}
