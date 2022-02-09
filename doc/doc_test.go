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

type testRowsReader byte

func (trr testRowsReader) Read([]byte) ([][]string, error) {
	return testRows, nil
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

	ed := NewRows(testRowsReader(0))
	err := ed.Encode([]byte{})
	assert.Assert(t, err == nil, err)
	// var ed rowsDoc
	// // ed.LoadFile()
	// ed.head = testRows[0]
	// ed.rows = testRows[1:]
	// rc := newRowsConv(testRows, reflect.TypeOf(*val))
	m := map[string]interface{}{"ID": val}
	err = ed.Decode(m) //rc.Unmarshal() //unmarshalRows(testRows, reflect.TypeOf(*val))
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
	js := Json(jsVal)
	var d Doc = &js

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

	assert.Equal(t, string(d.Val()), `{"Name":"hiank","Hope":"hope","Age":18}`)

	var tmpJs Json
	tmpJs.Encode(val2)
	assert.Equal(t, string(tmpJs.Val()), `{"Name":"hiank","Hope":"hope","Age":18}`)

	js2 := Json([]byte{})
	rv := reflect.ValueOf(js2)
	assert.Equal(t, rv.Kind(), reflect.Slice)

	js2.Encode(val2)
	assert.Equal(t, string(js2.Val()), `{"Name":"hiank","Hope":"hope","Age":18}`)
	assert.Equal(t, string((&js2).Val()), `{"Name":"hiank","Hope":"hope","Age":18}`)

	rv = reflect.ValueOf(Json([]byte{}))
	assert.Equal(t, rv.Kind(), reflect.Slice)

	rv = reflect.ValueOf(&js2)
	assert.Equal(t, rv.Kind(), reflect.Ptr)
}

func TestGob(t *testing.T) {
	var gb Gob
	d := &gb
	err := d.Encode(testStruct{Name: "gob", Hope: "ws"})
	assert.Assert(t, err == nil, err)

	var val2 testStruct
	d.Decode(&val2)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, val2.Name, "gob")
	assert.Equal(t, val2.Hope, "ws")
}

func TestBytesLenght(t *testing.T) {
	val := &testdata.Test1{Name: "hiank"}
	var js Json
	js.Encode(val)

	var pb PB
	pb.Encode(val)

	var gb Gob
	gb.Encode(val)

	var ym Yaml
	ym.Encode(val)

	assert.Assert(t, len(gb) > len(pb))
	assert.Assert(t, len(js) > len(ym), "jslen(%d) ymlen(%d)", len(js), len(ym))
	assert.Assert(t, len(ym) > len(pb), "pblen(%d) ymlen(%d)", len(pb), len(ym))
}

var testYamlStr = `name: host
m:
  Age: 11
  Lv: 22
  Id: 25`

type testYamlStruct struct {
	Name string
	M    map[string]int
}

func TestYaml(t *testing.T) {
	ym := Yaml([]byte(testYamlStr))
	// var val testYamlStruct2
	var val testYamlStruct //P: map[string]int{}}
	err := ym.Decode(&val)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, val.Name, "host")
	assert.Equal(t, val.M["Age"], 11)
	assert.Equal(t, val.M["Lv"], 22)
	assert.Equal(t, val.M["Id"], 25)

	// var outYm Yaml
	outYm := new(Yaml)
	err = outYm.Encode(val)
	assert.Assert(t, err == nil, err)
	var val2 testYamlStruct
	outYm.Decode(&val2)

	assert.Equal(t, val.Name, val2.Name)
	assert.Equal(t, val.M["Age"], val2.M["Age"])
	assert.Equal(t, val.M["Lv"], val2.M["Lv"])
	assert.Equal(t, val.M["Id"], val2.M["Id"])

	// assert.Equal(t, val.M, val2.M)
}
