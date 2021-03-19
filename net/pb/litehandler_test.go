package pb_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/net/pb/testdata"
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

func (th *testHandler) Handle(msg proto.Message) (err error) {
	select {
	case th.out <- msg:
	default:
		err = errors.New("testHandler out chan error")
	}
	return
}

func TestLiteHandler(t *testing.T) {

	outTest1 := make(chan proto.Message, 1)

	t.Run("HandleWithoutRegister", func(t *testing.T) {
		err := pb.LiteHandler.Handle(&testdata.Test1{Name: "noregister"})
		assert.Equal(t, err, codes.ErrorNoMessageHandler, "未注册任何Handler将返回错误")
	})

	t.Run("HandleWithoutRegisterButHasDefaultHandler", func(t *testing.T) {
		outDefault := make(chan proto.Message, 1)
		pb.LiteHandler.DefaultHandler = &testHandler{out: outDefault}
		err := pb.LiteHandler.Handle(&testdata.Test1{Name: "noregister"})
		assert.Assert(t, err == nil, "如果有DefaultHandler，则调用DefaultHandler")
		msg := <-outDefault
		assert.Equal(t, msg.(*testdata.Test1).GetName(), "noregister", "DefaultHandler将处理未匹配的消息")
		pb.LiteHandler.DefaultHandler = nil
	})

	t.Run("RegisterTest1", func(t *testing.T) {
		err := pb.LiteHandler.Register(new(testdata.Test1), &testHandler{out: outTest1})
		assert.Assert(t, err == nil)

		err = pb.LiteHandler.Register(new(testdata.Test1), &testHandler{out: outTest1})
		assert.Equal(t, err.Error(), codes.ErrorExistedMessageHandler.Error())
	})

	t.Run("RegisterNilHandler", func(t *testing.T) {

	})

	t.Run("HandleJustRegisterOtherMessage", func(t *testing.T) {
		err := pb.LiteHandler.Handle(&testdata.Test2{Hope: "ws"})
		assert.Equal(t, err, codes.ErrorNoMessageHandler, "未注册的消息，不会处理(无DefaultHandler)")
	})

	t.Run("HandleRegisterThenSetDefaultHandler", func(t *testing.T) {

		err := pb.LiteHandler.Handle(&testdata.Test1{Name: "alreadyregister"})
		assert.Assert(t, err == nil, "已注册的消息，必须能执行Handler")

		msg := <-outTest1
		assert.Equal(t, msg.(*testdata.Test1).GetName(), "alreadyregister")

		outDefault := make(chan proto.Message, 1)
		pb.LiteHandler.DefaultHandler = &testHandler{out: outDefault}

		err = pb.LiteHandler.Handle(&testdata.Test1{Name: "afterdefault"})
		assert.Assert(t, err == nil)

		err = pb.LiteHandler.Handle(&testdata.Test2{Hope: "hope"})
		assert.Assert(t, err == nil)

		msg1 := <-outTest1
		assert.Equal(t, msg1.(*testdata.Test1).Name, "afterdefault", "设置DefaultHandler不能影响注册Handler")

		msg2 := <-outDefault
		assert.Equal(t, msg2.(*testdata.Test2).Hope, "hope", "未注册是消息走DefaultHandler")

		pb.LiteHandler.DefaultHandler = nil
	})

}

func TestLiteHandlerHandleAny(t *testing.T) {

	outAnyTest1 := make(chan proto.Message, 1)
	outAny := make(chan proto.Message, 1)
	//
	t.Run("HandleEmptyAny", func(t *testing.T) {
		anyMsg := new(anypb.Any)
		err := pb.LiteHandler.Handle(anyMsg)
		assert.Equal(t, err, codes.ErrorAnyMessageIsEmpty, "空Any默认不会被处理")
	})

	t.Run("HandleEmptyAnyExistDefaultHandler", func(t *testing.T) {

		pb.LiteHandler.DefaultHandler = &testHandler{}
		anyMsg := new(anypb.Any)
		err := pb.LiteHandler.Handle(anyMsg)
		assert.Equal(t, err, codes.ErrorAnyMessageIsEmpty, "即使存在DefaultHandler的情况下，空Any也不会被调用")

		pb.LiteHandler.DefaultHandler = nil
	})

	t.Run("HandleRegisteredAnyMessage", func(t *testing.T) {

		pb.LiteHandler.Register(&testdata.AnyTest1{}, &testHandler{out: outAnyTest1})

		anyMsg, _ := anypb.New(&testdata.AnyTest1{Name: "test1"})
		err := pb.LiteHandler.Handle(anyMsg)
		assert.Assert(t, err == nil, "Any消息会被注册的Handler处理")

		msg := <-outAnyTest1
		assert.Assert(t, !reflect.DeepEqual(anyMsg, msg), "如果有注册的Handler，Any会自动解包")
		assert.Equal(t, msg.(*testdata.AnyTest1).GetName(), "test1", "Any消息会被自动解析为对应的消息")
	})

	t.Run("HandleAnyTest2", func(t *testing.T) {
		anyMsg, _ := anypb.New(&testdata.AnyTest2{Hope: "oh"})
		err := pb.LiteHandler.Handle(anyMsg)
		assert.Equal(t, err, codes.ErrorNoMessageHandler, "Any消息未注册的，不会被处理")
	})

	t.Run("HandleAnyTest2ExistDefaultHandler", func(t *testing.T) {

		outDefault := make(chan proto.Message, 1)
		pb.LiteHandler.DefaultHandler = &testHandler{out: outDefault}
		anyMsg, _ := anypb.New(&testdata.AnyTest2{Hope: "wow"})
		err := pb.LiteHandler.Handle(anyMsg)
		assert.Assert(t, err == nil, "DefaultHandler可以处理为注册的Any消息")

		msg := <-outDefault
		assert.Assert(t, reflect.DeepEqual(msg, anyMsg), "Any被DefaultHandler处理的情况下，不能解包，避免下一层使用者无法获知消息必要信息")

		pb.LiteHandler.DefaultHandler = nil
	})

	t.Run("HandleAnyWithRegister", func(t *testing.T) {

		pb.LiteHandler.Register(new(anypb.Any), &testHandler{out: outAny})

		anyMsg1, _ := anypb.New(&testdata.AnyTest1{Name: "any1"})
		err := pb.LiteHandler.Handle(anyMsg1)
		assert.Assert(t, err == nil)

		msg := <-outAny
		assert.Assert(t, reflect.DeepEqual(anyMsg1, msg), "注册Any消息Handler，直接调用此Handler，而不是使用映射对象的Handler")
	})

	t.Run("HandleAnyWithRegisterExistDefaultHandler", func(t *testing.T) {

		outDefault := make(chan proto.Message, 1)
		pb.LiteHandler.DefaultHandler = &testHandler{out: outDefault}

		pb.LiteHandler.Register(new(anypb.Any), &testHandler{out: outAny})

		anyMsg1, _ := anypb.New(&testdata.AnyTest1{Name: "any1"})
		err := pb.LiteHandler.Handle(anyMsg1)
		assert.Assert(t, err == nil)

		msg := <-outAny
		assert.Assert(t, reflect.DeepEqual(anyMsg1, msg), "注册Any消息Handler，直接调用此Handler，而不是使用DefaultHandler")

		pb.LiteHandler.DefaultHandler = nil
	})
}

func TestLiteHandlerHandlePBMessage(t *testing.T) {

	// outMessage := make(chan proto.Message, 1)
	outTest1, outPB := make(chan proto.Message, 1), make(chan proto.Message, 1)

	// pb.LiteHandler.Register(&testdata.MessageTest{}, &testHandler{out: outTest})
	// pb.LiteHandler.Register(&pb.Message{}, &testHandler{out: outMessage})

	t.Run("HandleNothingRegistered", func(t *testing.T) {
		anyMsg, _ := anypb.New(&testdata.MessageTest1{Key: "hello"})
		msg := &pb.Message{Value: anyMsg}
		err := pb.LiteHandler.Handle(msg)
		assert.Equal(t, err, codes.ErrorNoMessageHandler)
	})

	t.Run("HandleValueRegisterd", func(t *testing.T) {

		pb.LiteHandler.Register(&testdata.MessageTest1{}, &testHandler{out: outTest1})

		anyMsg, _ := anypb.New(&testdata.MessageTest1{Key: "world"})
		msg := &pb.Message{Value: anyMsg}
		err := pb.LiteHandler.Handle(msg)
		assert.Assert(t, err == nil, "针对Message消息，默认会找到Value的期望Handler，并处理")

		rltMsg := <-outTest1
		assert.Assert(t, reflect.DeepEqual(rltMsg, msg), "针对Message消息，不会只把解析出的Value传出")
	})

	t.Run("HandlePBMessageRegistered", func(t *testing.T) {

		pb.LiteHandler.Register(&pb.Message{}, &testHandler{out: outPB})

		anyMsg, _ := anypb.New(&testdata.MessageTest1{Key: "like"})
		msg := &pb.Message{Value: anyMsg}
		err := pb.LiteHandler.Handle(msg)
		assert.Assert(t, err == nil)

		select {
		case rlt := <-outPB:
			assert.Assert(t, reflect.DeepEqual(rlt, msg))
		case <-outTest1:
			assert.Assert(t, false, "如果注册了pb.Message的Handler，则优先使用之")
		}
	})

}
