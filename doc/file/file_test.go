package file_test

import (
	"strconv"
	"testing"

	"github.com/hiank/think/doc/file"
	"gotest.tools/v3/assert"
)

type testExcelConfig struct {
	ID   uint   `excel:"关卡ID"`
	Name string `excel:"关卡名"`
	Lv   uint   `excel:"怪物等级"`
	T1   uint   `excel:"队伍1神兽"`
}

var wantExcel = []*testExcelConfig{
	{ID: 1, Name: "唐朝精锐", Lv: 100, T1: 20303},
	{ID: 2, Name: "唐朝精锐", Lv: 100, T1: 20301},
	{ID: 3, Name: "唐朝精锐", Lv: 100, T1: 20303},
	{ID: 4, Name: "明朝精锐", Lv: 100, T1: 20301},
	{ID: 5, Name: "明朝精锐", Lv: 100, T1: 20303},
	{ID: 6, Name: "明朝精锐", Lv: 100, T1: 20303},
	{ID: 7, Name: "宋朝精锐", Lv: 100, T1: 20303},
	{ID: 8, Name: "宋朝精锐", Lv: 100, T1: 20303},
	{ID: 9, Name: "宋朝精锐", Lv: 100, T1: 20303},
	{ID: 10, Name: "宋朝精锐", Lv: 100, T1: 20303},
}

func TestFit(t *testing.T) {
	t.Run("FormRows", func(t *testing.T) {
		buffer := file.Fit(file.FormRows)
		// buffer.LoadBytes(file.FormRows, )
		err := buffer.LoadFile("testdata/config.xlsx")
		assert.Assert(t, err == nil, err)

		var val testExcelConfig
		err = buffer.Decode(&val)
		assert.Assert(t, err != nil, "only support array or map param")

		m := map[string]interface{}{}
		m["ID"] = new(testExcelConfig)
		err = buffer.Decode(m)
		assert.Assert(t, err == nil, err)

		for _, v := range wantExcel {
			assert.Equal(t, *m[strconv.FormatUint(uint64(v.ID), 10)].(*testExcelConfig), *v)
		}

		l := []interface{}{new(testExcelConfig)}
		err = buffer.Decode(&l)
		assert.Assert(t, err == nil, err)
		for i, v := range wantExcel {
			assert.Equal(t, *l[i].(*testExcelConfig), *v)
		}

		l = []interface{}{testExcelConfig{}}
		err = buffer.Decode(&l)
		assert.Assert(t, err == nil, err)
		for i, v := range wantExcel {
			assert.Equal(t, *l[i].(*testExcelConfig), *v)
		}
	})
	t.Run("Json", func(t *testing.T) {
		buffer := file.Fit(file.FormJson)
		buffer.LoadBytes(file.FormJson, []byte(`{"sys.Tik": "nil"}`))

		var cfg testConfig
		buffer.Decode(&cfg)
		assert.Equal(t, cfg.Tik, "nil")

		buffer.LoadBytes(file.FormJson, []byte(`{"sys.Tik": "overwrite"}`))
		buffer.Decode(&cfg)
		assert.Equal(t, cfg.Tik, "overwrite", "overwrite previous settings")
	})
	t.Run("Yaml", func(t *testing.T) {
		buffer := file.Fit(file.FormYaml)
		err := buffer.LoadBytes(file.FormGob, []byte(`Tok: 2`))
		assert.Assert(t, err != nil, "form mismatch")

		buffer.LoadBytes(file.FormYaml, []byte(`Tok: 2`))
		var cfg testConfig
		buffer.Decode(&cfg)

		assert.Equal(t, cfg.Tok, 2)

		buffer.LoadBytes(file.FormYaml, []byte(`Tok: 3`))
		buffer.Decode(&cfg)
		assert.Equal(t, cfg.Tok, 3, "overwrite previous settings")
	})
}

type testConfig struct {
	Tik string `json:"sys.Tik"`
	Tok int    `yaml:"Tok"`
}
