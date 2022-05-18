package box_test

import (
	"testing"

	"github.com/hiank/think/net/testdata"
	"github.com/hiank/think/net/box"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

func TestNew(t *testing.T) {
	msg, err := box.New(&testdata.AnyTest1{Name: "test new"})
	assert.Equal(t, err, nil, err)
	at1, err := msg.GetAny().UnmarshalNew()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, at1.(*testdata.AnyTest1).GetName(), "test new")

	amsg2, _ := anypb.New(&testdata.AnyTest2{Hope: "audi"})
	msg, err = box.New(amsg2)
	assert.Equal(t, err, nil, err)
	at2, err := msg.GetAny().UnmarshalNew()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, at2.(*testdata.AnyTest2).GetHope(), "audi")
}

func TestUnmarshal(t *testing.T) {
	b, _ := proto.Marshal(&testdata.AnyTest1{Name: "benz"})
	var msg box.Message
	err := box.Unmarshal[*testdata.AnyTest1](b, &msg)
	assert.Equal(t, err, nil, err)

	at1, err := msg.GetAny().UnmarshalNew()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, at1.(*testdata.AnyTest1).GetName(), "benz")

	m, err := box.UnmarshalNew[*testdata.AnyTest1](b)
	assert.Equal(t, err, nil, err)
	at1, err = m.GetAny().UnmarshalNew()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, at1.(*testdata.AnyTest1).GetName(), "benz")
}

func TestAnyNew(t *testing.T) {
	amsg, _ := anypb.New(&testdata.AnyTest1{Name: "cw"})
	v2, err := anypb.New(amsg)
	assert.Equal(t, err, nil, err)
	v3, err := v2.UnmarshalNew()
	assert.Equal(t, err, nil, err)
	amsg, ok := v3.(*anypb.Any)
	assert.Assert(t, ok)
	v4, _ := amsg.UnmarshalNew()
	assert.Equal(t, v4.(*testdata.AnyTest1).GetName(), "cw")
}
