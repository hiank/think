package pb_test

import (
	"testing"

	"github.com/hiank/think/net/pb"
	"github.com/hiank/think/net/testdata"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"gotest.tools/v3/assert"
)

func TestGetServeName(t *testing.T) {
	// testdata.G_Example
	val := &testdata.G_Example{Value: "hh"}
	anyMsg, err := anypb.New(val)
	assert.Assert(t, err == nil, err)
	svcName, err := pb.GetServeName(anyMsg)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, svcName, "example")
}

func TestGetServeType(t *testing.T) {
	msgs := []protoreflect.ProtoMessage{
		&testdata.G_Example{},
		&testdata.P_Example{},
		&testdata.S_Example{},
		&testdata.TEST_Example{},
		&testdata.AnyTest1{},
	}
	wantypes := []int{
		pb.TypeGET,
		pb.TypePOST,
		pb.TypeSTREAM,
		pb.TypeUndefined,
		pb.TypeUndefined,
	}

	for i, msg := range msgs {
		anyMsg, err := anypb.New(msg)
		assert.Assert(t, err == nil, err)
		tp, _ := pb.GetServeType(anyMsg)
		assert.Equal(t, tp, wantypes[i])
	}

	_, err := pb.GetServeType(&anypb.Any{})
	assert.Assert(t, err != nil)
}
