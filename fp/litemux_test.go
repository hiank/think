package fp_test

import (
	"reflect"
	"testing"

	"github.com/hiank/think/fp"
	"gotest.tools/v3/assert"
)

type testConfig struct {
	Limit int    `json:"sys.Limit"`
	Key   string `yaml:"a"`
}

func TestParser(t *testing.T) {
	// t.Run()
	u := fp.NewParser()
	u.LoadFile("testdata", "testdata/config.json")

	var cfg testConfig
	u.ParseAndClear(&cfg)

	assert.Equal(t, cfg.Key, "love-ws")
	assert.Equal(t, cfg.Limit, 2)
}

type testExcelConfig struct {
	ID   uint   `excel:"关卡ID"`
	Name string `excel:"关卡名"`
	Lv   uint   `excel:"怪物等级"`
	T1   uint   `excel:"队伍1神兽"`
}

func TestUnmarshalExcel(t *testing.T) {
	vals := fp.UnmarshalExcel("testdata/config.xlsx", reflect.TypeOf(testExcelConfig{}))
	wantVals := []*testExcelConfig{
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
	assert.Equal(t, len(vals), len(wantVals))

	for i, val := range vals {
		wantVal, v := wantVals[i], val.(*testExcelConfig)
		assert.Equal(t, wantVal.ID, v.ID)
		assert.Equal(t, wantVal.Name, v.Name)
		assert.Equal(t, wantVal.Lv, v.Lv)
		assert.Equal(t, wantVal.T1, v.T1)
	}
}
