package doc_test

import (
	"testing"

	"github.com/hiank/think/doc"
	"github.com/hiank/think/doc/excel"
	"github.com/hiank/think/doc/testdata"

	"gotest.tools/v3/assert"
)

type tmpJson struct {
	Id   int    `json:"tmp.id"`
	Name string `json:"tmp.name"`
}

type tmpYaml struct {
	Key   string `yaml:"ws.k"`
	Value int    `yaml:"ws.v"`
}

type tmpGob struct {
	Age  int
	Hope string
}

const tmpJsonValue = `{
	"tmp.id": 9527,
	"tmp.name": "华安"
}`
const tmpYamlValue = `
ws.k: "bbq"
ws.v: 25
`

func TestCoder(t *testing.T) {
	t.Run("JsonCoder", func(t *testing.T) {
		coder := doc.NewCoder[doc.JsonCoder]()
		val := tmpJson{}
		err := coder.Decode([]byte(tmpJsonValue), &val)
		assert.Equal(t, err, nil, err)
		assert.DeepEqual(t, val, tmpJson{Id: 9527, Name: "华安"})

		err = coder.Decode([]byte(tmpJsonValue), val)
		assert.Assert(t, err != nil, "non-pointer")

		b, err := coder.Encode(tmpJson{Id: 8821, Name: "love ws"})
		assert.Equal(t, err, nil, err)
		assert.Equal(t, string(b), `{"tmp.id":8821,"tmp.name":"love ws"}`)

		b, err = coder.Encode(&tmpJson{Id: 23, Name: "hiank"})
		assert.Equal(t, err, nil, err)
		assert.Equal(t, string(b), `{"tmp.id":23,"tmp.name":"hiank"}`)
	})

	t.Run("PBCoder", func(t *testing.T) {
		coder := doc.NewCoder[doc.PBCoder]()
		var msg = testdata.Test1{Name: "110"}
		// var msg testdata.Test1
		b, err := coder.Encode(&msg)
		assert.Equal(t, err, nil, err)

		var msg2 testdata.Test1
		err = coder.Decode(b, &msg2)
		assert.Equal(t, err, nil, err)
		assert.Equal(t, msg2.GetName(), "110")

		// _, err = coder.Encode(msg)
		// assert.Equal(t, err, doc.ErrNotProtoMessage)
	})

	t.Run("GobCoder", func(t *testing.T) {
		coder := doc.NewCoder[doc.GobCoder]()
		var v = tmpGob{Age: 18, Hope: "always"}
		b, err := coder.Encode(v)
		assert.Equal(t, err, nil, err)

		var v2 tmpGob
		err = coder.Decode(b, &v2)
		assert.Equal(t, err, nil, err)
		assert.DeepEqual(t, v, v2)
	})

	t.Run("YamlCoder", func(t *testing.T) {
		coder := doc.NewCoder[doc.YamlCoder]()
		// var v = tmpYaml{Key: "k", Value: 99}
		var v tmpYaml
		err := coder.Decode([]byte(tmpYamlValue), &v)
		assert.Equal(t, err, nil, err)
		assert.DeepEqual(t, v, tmpYaml{Key: "bbq", Value: 25})
		// coder.Encode(tmpYamlValue, )

		b, err := coder.Encode(tmpYaml{Key: "love ws", Value: 18})
		assert.Equal(t, err, nil, err)
		assert.Equal(t, string(b), "ws.k: love ws\nws.v: 18\n")
	})
}

func TestMakeT(t *testing.T) {
	v, err := doc.MakeT[tmpJson]()
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, v, tmpJson{Id: 0, Name: ""})

	pv, err := doc.MakeT[*tmpJson]()
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, pv, &tmpJson{Id: 0, Name: ""})

	_, err = doc.MakeT[**tmpJson]()
	assert.Equal(t, err, doc.ErrUnsupportType)

	iv, err := doc.MakeT[int]()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, iv, 0)

	// iv = 11
	piv, err := doc.MakeT[*int]()
	assert.Equal(t, err, nil, err)
	assert.Equal(t, *piv, 0)

	*piv = 11
	assert.Equal(t, *piv, 11)

	mv, err := doc.MakeT[map[string]string]()
	assert.Equal(t, err, nil, err)
	mv["hp"] = "low"
	assert.DeepEqual(t, mv, map[string]string{"hp": "low"})

	_, err = doc.MakeT[*map[string]string]()
	assert.Equal(t, err, doc.ErrUnsupportType)

	sv, err := doc.MakeT[[]string]()
	assert.Equal(t, err, nil, err)
	sv = append(sv, []string{"11", "12"}...)
	assert.Equal(t, len(sv), 2)
	assert.Equal(t, cap(sv), 2)

	sv, err = doc.MakeT[[]string](2, 4)
	assert.Equal(t, err, nil, err)
	assert.Equal(t, len(sv), 2)
	assert.Equal(t, cap(sv), 4)

	_, err = doc.MakeT[*[]string]()
	assert.Equal(t, err, doc.ErrUnsupportType)
}

func TestDoc(t *testing.T) {
	d := doc.New[tmpJson](doc.NewCoder[doc.JsonCoder]())
	v, err := d.DecodeNew([]byte(tmpJsonValue))
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, v, tmpJson{Id: 9527, Name: "华安"})

	d2 := doc.New[*tmpJson](doc.NewCoder[doc.JsonCoder]())
	var v2 tmpJson
	err = d2.Decode([]byte(tmpJsonValue), &v2)
	assert.Equal(t, err, nil)
	assert.Equal(t, d2.T(), &v2)
	assert.DeepEqual(t, v2, tmpJson{Id: 9527, Name: "华安"})
	assert.DeepEqual(t, d2.Bytes(), []byte(tmpJsonValue))

	b, err := d2.Encode(&tmpJson{Id: 121, Name: "hiank"})
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, b, d2.Bytes())
	assert.Equal(t, string(b), `{"tmp.id":121,"tmp.name":"hiank"}`)
	assert.DeepEqual(t, d2.T(), &tmpJson{Id: 121, Name: "hiank"})
}

type tmpExcel struct {
	Lv   uint   `excel:"怪物等级"`
	ID   string `excel:"关卡ID"`
	Name string `excel:"名字"`
}

var tmpRows [][]string = [][]string{
	{"关卡ID", "怪物等级", "关卡名字"},
	{"11", "12", "无知"},
	{"1", "2", "优质"},
}

type tmpRowDecoder struct{}

func (tmpRowDecoder) UnmarshalNew([]byte) (excel.Rows, error) {
	return tmpRows, nil
}

func TestRowsCoder(t *testing.T) {
	coder := doc.NewRowsCoder[tmpExcel]("ID", tmpRowDecoder{})
	// var m map[string]*tmpExcel
	err := coder.Decode([]byte{}, make(map[string]*tmpExcel))
	assert.Equal(t, err, doc.ErrInvalidParamType)

	var m map[string]tmpExcel = make(map[string]tmpExcel)
	err = coder.Decode([]byte{}, m)
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, m, map[string]tmpExcel{"11": {ID: "11", Lv: 12}, "1": {ID: "1", Lv: 2}})

	coder = doc.NewRowsCoder[tmpExcel]("Non", tmpRowDecoder{})
	err = coder.Decode([]byte{}, m)
	assert.Equal(t, err, excel.ErrNonKeyFound)

	// coder = doc.NewRowsCoder[tmpExcel]("", tmpRowDecoder{})
	// err = coder.Decode([]byte{}, m)
	// assert.Equal(t, err, doc.ErrInvalidParamType)

	var sv []tmpExcel
	err = coder.Decode([]byte{}, &sv)
	assert.Equal(t, err, nil, nil)
	assert.DeepEqual(t, sv, []tmpExcel{{ID: "11", Lv: 12}, {ID: "1", Lv: 2}})
}
