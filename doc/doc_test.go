package doc_test

import (
	"testing"

	"github.com/hiank/think/doc"
	"github.com/hiank/think/doc/testdata"
	"gotest.tools/v3/assert"
)

type testJson struct {
	Id   int
	Name string
	Age  int `json:"tag.age"`
	Non  bool
}

const testJsonStr = `{"id":112,"Name":"hiank","tag.age":31,"NON":true}`

type testYaml struct {
	Name string
	ID   int
	L    []int `yaml:"list"`
	M    map[string]int
}

var testYamlStr = `name: host
id: 25
list:
- 320
- 325
- 330
- 340
m:
  Id: 25
  Lv: 22
  age: 11
`

type testExcel struct {
	Lv   uint   `excel:"怪物等级"`
	ID   int    `excel:"关卡ID"`
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

func verifyT(tc doc.T, data []byte, final any, t *testing.T) {
	_, err := tc.Encode()
	assert.Equal(t, err, nil, "encode empty V")
	// assert.Equal(t, len(b), 0)

	err = tc.Decode(data)
	assert.Equal(t, err, nil)

	//
	assert.DeepEqual(t, tc.V, final)

	//
	b, err := tc.Encode()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(b), len(data))
}

func TestT(t *testing.T) {
	tc := doc.T{V: &testJson{}}
	err := tc.Decode([]byte(testJsonStr))
	assert.Assert(t, err != nil, "no coder for T. only support make by Maker")

	// val := &testJson{}
	var val testJson
	tc = doc.J.MakeT(&val)
	err = tc.Decode([]byte(testJsonStr))
	assert.Equal(t, err, nil)
	// json decode will ignore all character case
	assert.DeepEqual(t, val, testJson{Id: 112, Name: "hiank", Age: 31, Non: true})

	val.Id = 9
	val.Non = false
	b, err := tc.Encode()
	assert.Equal(t, err, nil)
	assert.Equal(t, string(b), `{"Id":9,"Name":"hiank","tag.age":31,"Non":false}`)

	tc = doc.J.MakeT(val)
	err = tc.Decode([]byte(testJsonStr))
	assert.Assert(t, err != nil, "T's V must be a pointer for struct")

	b2, err := tc.Encode()
	assert.Equal(t, err, nil, "encode is ok")
	assert.DeepEqual(t, b, b2)

	t.Run("json", func(t *testing.T) {
		var val testJson
		tc := doc.J.MakeT(val)
		err := tc.Decode([]byte(testJsonStr))
		assert.Assert(t, err != nil, "cannot decode to not pointer v")

		target := &testJson{
			Id:   112,
			Name: "hiank",
			Age:  31,
			Non:  true,
		}
		verifyT(doc.J.MakeT(&testJson{}), []byte(testJsonStr), target, t)
	})

	t.Run("yaml", func(t *testing.T) {
		var val testYaml
		tc := doc.Y.MakeT(val)
		err := tc.Decode([]byte(testYamlStr))
		assert.Assert(t, err != nil, "cannot decode to not pointer v")

		target := &testYaml{
			ID:   25,
			Name: "host",
			L:    []int{320, 325, 330, 340},
			M: map[string]int{
				"age": 11,
				"Lv":  22,
				"Id":  25,
			},
		}
		verifyT(doc.Y.MakeT(&testYaml{}), []byte(testYamlStr), target, t)
	})

	t.Run("gob", func(t *testing.T) {
		target := &testYaml{
			Name: "gob",
			ID:   1,
			L:    []int{110, 111},
			M: map[string]int{
				"hh": 2,
				"ip": 10,
			},
		}
		// var val testYaml
		tc := doc.G.MakeT(target)

		b, err := tc.Encode()
		assert.Equal(t, err, nil)

		verifyT(doc.G.MakeT(&testYaml{}), b, target, t)
	})

	t.Run("protobuf", func(t *testing.T) {
		target := &testdata.Test1{Name: "test protobuf T"}
		// var val testYaml
		tc := doc.P.MakeT(target)

		b, err := tc.Encode()
		assert.Equal(t, err, nil)

		var val testdata.Test1
		tc = doc.P.MakeT(&val)
		tc.Decode(b)
		assert.Equal(t, val.Name, target.Name)
	})

	t.Run("rows", func(t *testing.T) {
		maker := doc.NewMaker(&doc.RowsCoder{RC: testRowsConverter(0)})
		target := []*testExcel{}
		tc := maker.MakeT(&target)

		//RC will convert any data to testRows
		err := tc.Decode([]byte{})
		assert.Equal(t, err, nil)
		assert.Equal(t, len(target), 2)

		tm := map[int]*testExcel{}
		maker = doc.NewMaker(&doc.RowsCoder{KT: "ID", RC: testRowsConverter(0)})

		tc = maker.MakeT(tm)
		err = tc.Decode([]byte{})
		assert.Equal(t, err, nil)

		assert.Equal(t, len(tm), 2)
		assert.DeepEqual(t, tm[11], &testExcel{
			ID: 11,
			Lv: 12,
		})
		assert.DeepEqual(t, tm[1], &testExcel{
			ID: 1,
			Lv: 2,
		})
	})
}

func TestB(t *testing.T) {
	bc := &doc.B{D: []byte(testJsonStr)}
	err := bc.Decode(&testJson{})
	assert.Assert(t, err != nil, "no coder for B. only support make by Maker")

	// val := &testJson{}
	// var val testJson

	t.Run("json", func(t *testing.T) {
		bc := doc.J.MakeB(nil)
		err := bc.Encode([]byte(testJsonStr))
		assert.Equal(t, err, nil)
		assert.DeepEqual(t, bc.D, []byte(testJsonStr))

		var val testJson
		bc.Decode(&val)
		// json decode will ignore all character case
		assert.DeepEqual(t, val, testJson{Id: 112, Name: "hiank", Age: 31, Non: true})
	})

	t.Run("yaml", func(t *testing.T) {
		bc := doc.Y.MakeB([]byte(testYamlStr))
		assert.DeepEqual(t, bc.D, []byte(testYamlStr))

		var val testYaml
		bc.Decode(&val)
		// json decode will ignore all character case
		assert.DeepEqual(t, val, testYaml{
			ID:   25,
			Name: "host",
			L:    []int{320, 325, 330, 340},
			M: map[string]int{
				"age": 11,
				"Lv":  22,
				"Id":  25,
			},
		})

		val.M = map[string]int{"new": 2}
		bc.Encode(val)
		str1 := string(bc.D)
		bc.Encode(&val)
		str2 := string(bc.D)
		assert.Equal(t, str1, str2)

		var val2 testYaml
		bc.Decode(&val2)
		assert.DeepEqual(t, val, val2)
	})

	t.Run("gob", func(t *testing.T) {
		bc := doc.G.MakeB(nil)
		// assert.DeepEqual(t, bc.D, []byte(testYamlStr))
		val1 := testYaml{ID: 21, Name: "ws", L: []int{11, 12, 13}, M: map[string]int{"hope": 21}}
		bc.Encode(val1)
		var val2 testYaml
		bc.Decode(&val2)
		// json decode will ignore all character case
		assert.DeepEqual(t, val1, val2)
	})

	t.Run("protobuf", func(t *testing.T) {
		bc := doc.G.MakeB(nil)
		// assert.DeepEqual(t, bc.D, []byte(testYamlStr))
		val1 := &testdata.Test2{Age: 18}
		bc.Encode(val1)
		var val2 testdata.Test2
		bc.Decode(&val2)
		// json decode will ignore all character case
		assert.Equal(t, val1.Age, val2.Age)
	})

	t.Run("rows", func(t *testing.T) {
		maker := doc.NewMaker(&doc.RowsCoder{RC: testRowsConverter(0)})
		// target := []*testExcel{}
		bc := maker.MakeB(nil)

		l := []*testExcel{}
		//RC will convert any data to testRows
		err := bc.Decode(&l)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(l), 2)

		tm := map[int]*testExcel{}
		maker = doc.NewMaker(&doc.RowsCoder{KT: "ID", RC: testRowsConverter(0)})

		bc = maker.MakeB(nil)
		err = bc.Decode(tm)
		assert.Equal(t, err, nil)

		assert.Equal(t, len(tm), 2)
		assert.DeepEqual(t, tm[11], &testExcel{
			ID: 11,
			Lv: 12,
		})
		assert.DeepEqual(t, tm[1], &testExcel{
			ID: 1,
			Lv: 2,
		})
	})

}

func TestTcoder(t *testing.T) {
	var tcoder doc.Tcoder

	jv := &testJson{
		Id:   12,
		Name: "hp",
		Age:  39,
		Non:  true,
	}
	tc := doc.J.MakeT(jv)
	b, err := tcoder.Encode(tc)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(b), `{"Id":12,"Name":"hp","tag.age":39,"Non":true}`)

	err = tcoder.Decode([]byte(`{"Id":2,"Name":"hope","tag.age":40,"Non":false}`), tc)
	assert.Equal(t, err, nil)

	assert.DeepEqual(t, jv, &testJson{
		Id:   2,
		Name: "hope",
		Age:  40,
	})
}
