package doc

import (
	"reflect"
	"testing"

	"github.com/hiank/think/doc/testdata"
	"google.golang.org/protobuf/proto"
	"gotest.tools/v3/assert"
)

type testExcel struct {
	Lv   uint   `excel:"怪物等级"`
	ID   string `excel:"关卡ID"`
	Name string `excel:"关卡名字"`
}

var testRows [][]string = [][]string{
	{"关卡ID", "怪物等级", "名字"},
	{"11", "12", "无知"},
	{"1", "2", "优质"},
}

func TestRowsDoc(t *testing.T) {
	val := &testExcel{}
	fv := reflect.ValueOf(val)
	assert.Equal(t, fv.Kind(), reflect.Ptr)

	fv = fv.Elem()
	assert.Equal(t, fv.Kind(), reflect.Struct)

	ff, _ := fv.Type().FieldByName("ID")
	tag := ff.Tag.Get("excel")
	assert.Equal(t, tag, "关卡ID")
	assert.Equal(t, ff.Name, "ID")

	var ed rowsDoc
	// ed.LoadFile()
	ed.head = testRows[0]
	ed.rows = testRows[1:]
	// rc := newRowsConv(testRows, reflect.TypeOf(*val))
	m := map[string]interface{}{}
	m["ID"] = val
	err := ed.Decode(m) //rc.Unmarshal() //unmarshalRows(testRows, reflect.TypeOf(*val))
	assert.Assert(t, err == nil, err)
	assert.Equal(t, len(m), 2)
	assert.Equal(t, m["11"].(*testExcel).Lv, uint(12))
	assert.Equal(t, m["1"].(*testExcel).Lv, uint(2))

	l := []interface{}{val}
	err = ed.Decode(&l)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, len(l), 2)
	assert.Equal(t, l[0].(*testExcel).Lv, uint(12))
}

func TestPB(t *testing.T) {
	// var doc docPB
	msg := &testdata.Test1{Name: "ll"}
	buf, err := proto.Marshal(msg)
	assert.Assert(t, err == nil, err)

	d := PB(buf)

	var msg1 testdata.Test1
	err = d.Decode(&msg1)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, msg1.GetName(), "ll")

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
	d := Json(jsVal) //doc.JsonMaker.Make([]byte(jsVal))

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
	var gb Gob
	err := gb.Encode(testStruct{Name: "gob", Hope: "ws"})
	assert.Assert(t, err == nil, err)

	var val2 testStruct
	gb.Decode(&val2)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, val2.Name, "gob")
	assert.Equal(t, val2.Hope, "ws")
}

func TestBytesMaker(t *testing.T) {
	_, ok := PBMaker.Make(nil).(*PB)
	assert.Assert(t, ok)

	_, ok = YamlMaker.Make(nil).(*Yaml)
	assert.Assert(t, ok)

	_, ok = JsonMaker.Make(nil).(*Json)
	assert.Assert(t, ok)

	_, ok = GobMaker.Make(nil).(*Gob)
	assert.Assert(t, ok)

	_, ok = NewRows(nil).(*rowsDoc)
	assert.Assert(t, ok)
}
