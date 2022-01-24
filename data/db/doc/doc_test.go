package doc_test

import (
	"testing"

	"github.com/hiank/think/data/db/doc"
	"github.com/hiank/think/data/db/doc/testdata"
	"google.golang.org/protobuf/proto"
	"gotest.tools/v3/assert"
)

func TestPB(t *testing.T) {
	// var doc docPB
	msg := &testdata.Test1{Name: "ll"}
	buf, err := proto.Marshal(msg)
	assert.Assert(t, err == nil, err)

	// d := doc.PB(buf)
	d := doc.PBMaker.Make(buf)

	var msg1 testdata.Test1
	err = d.Decode(&msg1)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, msg1.GetName(), "ll")
	assert.Equal(t, string(buf), d.Val())

	// docStr := string(doc)
	d.Encode(&testdata.Test2{Age: 18})
	// assert.Equal(t, docStr, string(doc))

	var msg2 testdata.Test2
	d.Decode(&msg2)
	assert.Equal(t, msg2.GetAge(), int32(18))

	// assert.Equal(t, string(buf), d.Val())
}

type testStruct struct {
	Name string
	Hope string
	Age  int
}

func TestJson(t *testing.T) {
	jsVal := `{"Name": "ll", "age": 18}`
	d := doc.JsonMaker.Make([]byte(jsVal))

	var val testStruct
	err := d.Decode(&val)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, val.Name, "ll")

	val.Name = "hiank"
	val.Hope = "hope"
	err = d.Encode(&val)
	assert.Assert(t, err == nil, err)
	// assert.Equal(t, d.Val(), `{"Name":"hiank","Hope":"hope"}`)

	var val2 testStruct
	err = d.Decode(&val2)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, val2.Name, "hiank")
	assert.Equal(t, val2.Hope, "hope")
	assert.Equal(t, val2.Age, 18)
}

func TestGob(t *testing.T) {
	var gb doc.Gob
	err := gb.Encode(testStruct{Name: "gob", Hope: "ws"})
	assert.Assert(t, err == nil, err)

	var val2 testStruct
	gb.Decode(&val2)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, val2.Name, "gob")
	assert.Equal(t, val2.Hope, "ws")

	assert.Equal(t, string(gb), gb.Val())

}
