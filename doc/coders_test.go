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

type testRowsConverter byte

func (trr testRowsConverter) ToRows([]byte) ([][]string, error) {
	return testRows, nil
}

func (trr testRowsConverter) ToBytes([][]string) ([]byte, error) {
	return []byte{}, nil
}

func TestReflectNew(t *testing.T) {
	// v := &testExcel{}
	var v testExcel
	rt := reflect.TypeOf(v)
	rv := reflect.New(rt)
	assert.Equal(t, rv.Kind(), reflect.Ptr)
}

func TestRowsCoder(t *testing.T) {
	t.Run("reflect map slice", func(t *testing.T) {
		m, l := map[int]int{}, []string{}

		assert.Equal(t, reflect.TypeOf(m).Kind(), reflect.Map)
		assert.Equal(t, reflect.TypeOf(l).Kind(), reflect.Slice)
		assert.Equal(t, reflect.TypeOf(&l).Kind(), reflect.Ptr)
		// assert.Equal(t, reflect.TypeOf(arr).Kind(), reflect.Array)

		assert.Equal(t, reflect.TypeOf(m).Elem().Kind(), reflect.Int)
		assert.Equal(t, reflect.TypeOf(l).Elem().Kind(), reflect.String)
		// assert.Equal(t, reflect.TypeOf(&l).Elem().Kind(), reflect.String)

		assert.Equal(t, reflect.ValueOf(m).Kind(), reflect.TypeOf(m).Kind())

		mptr := map[int]*testExcel{}
		rt := reflect.TypeOf(mptr).Elem()
		assert.Equal(t, rt.Kind(), reflect.Ptr)
		assert.Equal(t, rt.Elem().Kind(), reflect.Struct)

		m2 := map[int]testExcel{}
		rt = reflect.TypeOf(m2).Elem()
		assert.Equal(t, rt.Kind(), reflect.Struct)
	})
	t.Run("refelct slice set", func(t *testing.T) {
		l := []string{"0"}
		rv := reflect.ValueOf(&l)
		// rv = rv.Elem()
		nrv := reflect.Append(rv.Elem(), reflect.ValueOf("1"), reflect.ValueOf("2"))
		rv.Elem().Set(nrv)

		assert.Equal(t, len(l), 3)
		assert.DeepEqual(t, l, []string{"0", "1", "2"})
		// reflect.MakeSlice(reflect.TypeOf(l), 0, 4)
		// assert.Equal(t, cap(l), 3)

		var out []string = []string{}
		setFunc := func(arr *[]string) {
			rv := reflect.ValueOf(arr).Elem()
			nrv := reflect.Append(rv, reflect.ValueOf("1"))
			rv.Set(nrv)
		}
		setFunc(&out)
		assert.Equal(t, len(out), 1)
		assert.Equal(t, out[0], "1")
	})

	t.Run("parseHead", func(t *testing.T) {
		rd := &RowsCoder{KT: "ID"}

		m, i := rd.parseHead(testRows[0], reflect.TypeOf(testExcel{}))
		assert.Equal(t, len(m), 2, "'名字'无对应字段")
		assert.Equal(t, i, 1)
	})

	t.Run("toMap", func(t *testing.T) {
		rd := &RowsCoder{KT: "ID"}
		m := make(map[string]*testExcel)
		err := rd.toMap(testRows, reflect.ValueOf(m))
		assert.Equal(t, err, nil, err)

		assert.Equal(t, len(m), len(testRows)-1)
		assert.DeepEqual(t, m["11"], &testExcel{ID: "11", Lv: 12, Name: ""})

		m2 := make(map[string]testExcel)
		err = rd.toMap(testRows, reflect.ValueOf(m2))
		assert.Equal(t, err, nil, err)

		assert.Equal(t, len(m2), len(testRows)-1)
		assert.DeepEqual(t, m2["1"], testExcel{ID: "1", Lv: 2})
	})

	t.Run("toArray", func(t *testing.T) {
		rd := &RowsCoder{}
		l := &[]*testExcel{}
		err := rd.toArray(testRows, reflect.ValueOf(l).Elem())
		assert.Equal(t, err, nil, err)

		assert.Equal(t, len(*l), len(testRows)-1)
		assert.DeepEqual(t, (*l)[0], &testExcel{ID: "11", Lv: 12})

		l2 := &[]testExcel{}
		err = rd.toArray(testRows, reflect.ValueOf(l2).Elem())
		assert.Equal(t, err, nil, err)

		assert.Equal(t, (*l2)[1], testExcel{ID: "1", Lv: 2})
	})
	t.Run("Decode-Encode", func(t *testing.T) {
		rd := &RowsCoder{KT: "ID", RC: testRowsConverter(0)}
		m, l := make(map[string]*testExcel), []*testExcel{}
		err := rd.Decode([]byte{}, m)
		assert.Equal(t, err, nil, err)
		err = rd.Decode([]byte{}, &l)
		assert.Equal(t, err, nil, err)

		assert.Equal(t, len(m), len(l))
		assert.Equal(t, len(m), 2)

		var invalid testExcel
		assert.Equal(t, rd.Decode([]byte{}, &invalid), ErrNotSliceptrOrMap)
		assert.Equal(t, rd.Decode([]byte{}, invalid), ErrNotSliceptrOrMap)

		_, err = rd.Encode(m)
		assert.Equal(t, err, ErrUnimplemented)

		rows := [][]string{}
		rd.Decode([]byte{}, &rows)
		assert.DeepEqual(t, rows, testRows)
	})
}

func TestProtoCoder(t *testing.T) {
	// var doc docPB
	msg := &testdata.Test1{Name: "ll"}
	buf, err := proto.Marshal(msg)
	assert.Assert(t, err == nil, err)

	coder := protoCoder{}
	// d := PB(buf)

	var msg1 testdata.Test1
	err = coder.Decode(buf, &msg1)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, msg1.GetName(), "ll")

	// docStr := string(doc)
	buf, err = coder.Encode(&testdata.Test2{Age: 18})
	assert.Equal(t, err, nil, err)

	var msg2 testdata.Test2
	coder.Decode(buf, &msg2)
	assert.Equal(t, msg2.GetAge(), int32(18))

	// assert.Equal(t, string(buf), d.Val())
}

type testStruct struct {
	Name string
	Hope string
	Age  int
}

func TestJsonCoder(t *testing.T) {
	jsVal := `{"Name": "ll", "age": 18}`
	// js := Json(jsVal)
	// var d Doc = &js
	coder := jsonCoder{}

	var val testStruct
	err := coder.Decode([]byte(jsVal), &val)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, val.Name, "ll")

	val.Name = "hiank"
	val.Hope = "hope"
	buf, err := coder.Encode(&val)
	assert.Assert(t, err == nil, err)
	// assert.Equal(t, d.Val(), `{"Name":"hiank","Hope":"hope"}`)

	var val2 testStruct
	err = coder.Decode(buf, &val2)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, val2.Name, "hiank")
	assert.Equal(t, val2.Hope, "hope")
	assert.Equal(t, val2.Age, 18)

	assert.Equal(t, string(buf), `{"Name":"hiank","Hope":"hope","Age":18}`)

	// var tmpJs Json
	buf, _ = coder.Encode(val2)
	assert.Equal(t, string(buf), `{"Name":"hiank","Hope":"hope","Age":18}`)

	// js2 := Json([]byte{})
	// rv := reflect.ValueOf(js2)
	// assert.Equal(t, rv.Kind(), reflect.Slice)

	// js2.Encode(val2)
	// assert.Equal(t, string(js2.Val()), `{"Name":"hiank","Hope":"hope","Age":18}`)
	// assert.Equal(t, string((&js2).Val()), `{"Name":"hiank","Hope":"hope","Age":18}`)

	// rv = reflect.ValueOf(Json([]byte{}))
	// assert.Equal(t, rv.Kind(), reflect.Slice)

	// rv = reflect.ValueOf(&js2)
	// assert.Equal(t, rv.Kind(), reflect.Ptr)
}

func TestGob(t *testing.T) {
	// var gb Gob
	// d := &gb
	var coder gobCoder
	buf, err := coder.Encode(testStruct{Name: "gob", Hope: "ws"})
	assert.Assert(t, err == nil, err)

	var val2 testStruct
	coder.Decode(buf, &val2)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, val2.Name, "gob")
	assert.Equal(t, val2.Hope, "ws")
}

func TestBytesLenght(t *testing.T) {
	val := &testdata.Test1{Name: "hiank"}
	// var js Json
	js, _ := jsonCoder{}.Encode(val)

	// var pb PB
	pb, _ := protoCoder{}.Encode(val)

	// var gb Gob
	gb, _ := gobCoder{}.Encode(val)

	// var ym Yaml
	ym, _ := yamlCoder{}.Encode(val)

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
	var coder yamlCoder
	// ym := Yaml([]byte(testYamlStr))
	// var val testYamlStruct2
	var val testYamlStruct //P: map[string]int{}}
	err := coder.Decode([]byte(testYamlStr), &val)
	assert.Assert(t, err == nil, err)
	assert.Equal(t, val.Name, "host")
	assert.Equal(t, val.M["Age"], 11)
	assert.Equal(t, val.M["Lv"], 22)
	assert.Equal(t, val.M["Id"], 25)

	// var outYm Yaml
	// outYm := new(Yaml)
	buf, err := coder.Encode(val)
	assert.Assert(t, err == nil, err)
	var val2 testYamlStruct
	coder.Decode(buf, &val2)

	assert.Equal(t, val.Name, val2.Name)
	assert.Equal(t, val.M["Age"], val2.M["Age"])
	assert.Equal(t, val.M["Lv"], val2.M["Lv"])
	assert.Equal(t, val.M["Id"], val2.M["Id"])

	// assert.Equal(t, val.M, val2.M)
}

type testJsonStruct struct {
	Name string
}

func TestTcoder(t *testing.T) {
	var coder Tcoder
	_, err := coder.Encode(&testJsonStruct{Name: "name"})
	assert.Equal(t, err, ErrValueMustBeT)

	_, err = coder.Encode(Y.MakeT(nil))
	assert.Equal(t, err, ErrNilValue)

	var val testJsonStruct
	err = coder.Decode(nil, &val)
	assert.Equal(t, err, ErrValueMustBeT)

	err = coder.Decode(nil, Y.MakeT(nil))
	assert.Equal(t, err, ErrNilValue)

	t.Run("json", func(t *testing.T) {
		buf, err := coder.Encode(J.MakeT(&testJsonStruct{Name: "hiank"}))
		assert.Assert(t, err == nil, err)
		assert.Equal(t, string(buf), `{"Name":"hiank"}`)

		var val testJsonStruct
		err = coder.Decode(nil, J.MakeT(&val))
		assert.Assert(t, err != nil)
		assert.Equal(t, val.Name, "")

		err = coder.Decode(buf, J.MakeT(&val)) //db.T{D: doc.JsonMaker.Make(), V: &val})
		assert.Equal(t, err, nil)
		assert.Equal(t, val.Name, "hiank")
	})
	t.Run("proto", func(t *testing.T) {
		buf, err := coder.Encode(P.MakeT(&testdata.Test1{Name: "pb"}))
		assert.Assert(t, err == nil, err)

		var pb2 = new(testdata.Test1)
		err = coder.Decode(nil, P.MakeT(pb2))
		assert.Equal(t, err, nil)
		assert.Equal(t, pb2.GetName(), "")

		err = coder.Decode(buf, P.MakeT(pb2))
		assert.Assert(t, err == nil, err)
		assert.Equal(t, pb2.GetName(), "pb")
	})
	t.Run("gob", func(t *testing.T) {
		buf, err := coder.Encode(G.MakeT(&testdata.Test1{Name: "gob"}))
		assert.Assert(t, err == nil, err)

		var pb2 = new(testdata.Test1)
		err = coder.Decode(nil, G.MakeT(pb2))
		assert.Assert(t, err != nil)
		assert.Equal(t, pb2.GetName(), "")

		err = coder.Decode(buf, G.MakeT(pb2))
		assert.Assert(t, err == nil, err)
		assert.Equal(t, pb2.GetName(), "gob")

	})
	t.Run("yaml", func(t *testing.T) {
		// buf, err := coder.Encode(db.T{D: doc.GobMaker.Make(), V: &})
		// var v string
		// var d doc.T = doc.Y{V: doc.V{v}}
		// d.Decode(nil, nil)

		// d = doc.Y{V: doc.V{v}}
		buf, err := coder.Encode(Y.MakeT(&testJsonStruct{Name: "hiank"}))
		assert.Assert(t, err == nil, err)
		assert.Equal(t, string(buf), "name: hiank\n")

		var val testJsonStruct
		err = coder.Decode(nil, Y.MakeT(&val))
		assert.Assert(t, err == nil, "nil data for yaml decode is allowed")
		assert.Equal(t, val.Name, "")

		err = coder.Decode(buf, Y.MakeT(&val)) //db.T{D: doc.JsonMaker.Make(), V: &val})
		assert.Equal(t, err, nil)
		assert.Equal(t, val.Name, "hiank")
	})
}

func TestT(t *testing.T) {
	tc := T{}
	b, err := tc.Encode()
	assert.Assert(t, err != nil)
	assert.Equal(t, len(b), 0)

	// tc.c = gobCoder{}
	// _, err = tc.Encode()
	// assert.Equal(t, err, ErrNilValue)
}

func TestB(t *testing.T) {
	b := &B{}
	err := b.Encode(&testStruct{})
	assert.Assert(t, err != nil)

	// var val testStruct
	// err = b.Encode(&val)

}

type testApi interface {
	Val() string
}

type testBase struct {
}

func (testBase) Val() string {
	return "base"
}

type testFat struct {
	testBase
}

func (testFat) Val() string {
	return "fat"
}

func TestMore(t *testing.T) {
	tf := new(testFat)
	var api testApi = tf
	assert.Equal(t, api.Val(), "fat")
}
