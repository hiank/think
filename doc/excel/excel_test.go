package excel_test

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/hiank/think/doc/excel"
	"gotest.tools/v3/assert"
)

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

func TestHeader(t *testing.T) {
	header, err := excel.NewHeader[tmpExcel](tmpRows[0])
	assert.Equal(t, err, nil, err)
	v, err := header.NewT(tmpRows[1])
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, v, tmpExcel{12, "11", ""})

	_, err = header.NewT([]string{"i1", "tt", "pige"})
	assert.Assert(t, err != nil)

	v, err = header.NewT([]string{"11"})
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, v, tmpExcel{ID: "11"})

	v, err = header.NewT([]string{"11", "12", "pp", "tag"})
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, v, tmpExcel{Lv: 12, ID: "11"})

	kidx := header.KeyIndex("ID")
	assert.Assert(t, kidx != -1)
	assert.Equal(t, tmpRows[1][kidx], "11")

	kidx = header.KeyIndex("Non")
	assert.Equal(t, kidx, -1)

	kidx = header.KeyIndex("Lv")
	assert.Assert(t, kidx != -1)
	assert.Equal(t, tmpRows[1][kidx], "12")

	kidx = header.KeyIndex("LV")
	assert.Equal(t, kidx, -1)
	assert.Equal(t, tmpRows[1][header.KeyIndex("Lv")], "12")

	_, err = header.NewT([]string{"11", "99999999999999999999"})
	assert.Assert(t, err != nil, "overflow uint64 value")

	header2, err := excel.NewHeader[*tmpExcel](tmpRows[0])
	assert.Equal(t, err, nil, err)
	v2, err := header2.NewT(tmpRows[1])
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, v2, &tmpExcel{ID: "11", Lv: 12})

	_, err = excel.NewHeader[**tmpExcel](tmpRows[0])
	assert.Equal(t, err, excel.ErrNotStructype)

	_, err = excel.NewHeader[int](tmpRows[0])
	assert.Equal(t, err, excel.ErrNotStructype)

	_, err = excel.NewHeader[map[string]tmpExcel](tmpRows[0])
	assert.Equal(t, err, excel.ErrNotStructype)

}

// func TestFieldByName(t *testing.T) {
// 	var rt = reflect.TypeOf(tmpExcel{})
// 	assert.Assert(t, ok)
// }

func TestRangeRows(t *testing.T) {
	err := excel.Export_rangeRows(tmpRows, func(h *excel.Header[tmpExcel], s []string) error {
		return io.EOF
	})
	assert.Equal(t, err, io.EOF)

	cnt := 0
	err = excel.Export_rangeRows(tmpRows, func(h *excel.Header[*tmpExcel], s []string) error {
		cnt++
		return nil
	})
	assert.Equal(t, err, nil, err)
	assert.Equal(t, cnt, len(tmpRows)-1)

	cnt = 0
	err = excel.Export_rangeRows([][]string{tmpRows[0]}, func(h *excel.Header[*tmpExcel], s []string) error {
		cnt++
		return nil
	})
	assert.Equal(t, err, excel.ErrInvalidParamType)
	assert.Equal(t, cnt, 0)

	err = excel.Export_rangeRows([][]string{}, func(h *excel.Header[*tmpExcel], s []string) error {
		cnt++
		return nil
	})
	assert.Equal(t, err, excel.ErrInvalidParamType)
	assert.Equal(t, cnt, 0)
}

func TestUnmarshalMap(t *testing.T) {
	// m, err := excel.UnmarshalNewMap[tmpExcel](tmpRows, "Lv")
	// assert.Equal(t, err, nil, err)
	// assert.DeepEqual(t, m, map[string]tmpExcel{"12": {"11", 12}, "2": {"1", 2}})

	m2 := make(map[string]*tmpExcel)
	err := excel.UnmarshaltoMap(tmpRows, m2, "ID")
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, m2, map[string]*tmpExcel{"11": {ID: "11", Lv: 12}, "1": {ID: "1", Lv: 2}})
}

func TestUnmarshalSlice(t *testing.T) {
	// s, err := excel.UnmarshalNewSlice[*tmpExcel](tmpRows)
	// assert.Equal(t, err, nil, err)
	// assert.DeepEqual(t, s, []*tmpExcel{{"11", 12}, {"1", 2}})

	var out []tmpExcel
	err := excel.UnmarshaltoSlice(tmpRows, &out)
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, out, []tmpExcel{{ID: "11", Lv: 12}, {ID: "1", Lv: 2}})
}

// var tmpConfigRows = [][]string{
// 	{"1", "唐朝精锐", "100", "20303", "11378;11379;11380;11381;11382", "20304", "11378;11379;11380;11381;11382"},
// 	{"2", "唐朝精锐", "100", "20301", "11378;11379;11380;11381;11382", "20302", "11378;11379;11380;11381;11382"},
// 	{"3", "唐朝精锐", "100", "20303", "11378;11379;11380;11381;11382", "20304", "11378;11379;11380;11381;11382"},
// 	{"4", "明朝精锐", "100", "20301", "11378;11379;11380;11381;11382", "20302", "11378;11379;11380;11381;11382"},
// 	{"5", "明朝精锐", "100", "20303", "11378;11379;11380;11381;11382", "20304", "11378;11379;11380;11381;11382"},
// 	{"6", "明朝精锐", "100", "20303", "11378;11379;11380;11381;11382", "20304", "11378;11379;11380;11381;11382"},
// 	{"7", "宋朝精锐", "100", "20303", "11378;11379;11380;11381;11382", "20304", "11378;11379;11380;11381;11382"},
// 	{"8", "宋朝精锐", "100", "20303", "11378;11379;11380;11381;11382", "20304", "11378;11379;11380;11381;11382"},
// 	{"9", "宋朝精锐", "100", "20303", "11378;11379;11380;11381;11382", "20304", "11378;11379;11380;11381;11382"},
// 	{"10", "宋朝精锐", "100", "20303", "11378;11379;11380;11381;11382"},
// }

type tmpExcelConfig struct {
	ID   uint   `excel:"关卡ID"`
	Name string `excel:"关卡名"`
	Lv   uint   `excel:"怪物等级"`
	T1   uint   `excel:"队伍1神兽"`
	TM1  string `excel:"队伍1"`
	T2   uint   `excel:"队伍2神兽"`
	TM2  string `excel:"队伍2"`
}

var wantExcelConfig = []*tmpExcelConfig{
	{ID: 1, Name: "唐朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
	{ID: 2, Name: "唐朝精锐", Lv: 100, T1: 20301, TM1: "11378;11379;11380;11381;11382", T2: 20302, TM2: "11378;11379;11380;11381;11382"},
	{ID: 3, Name: "唐朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
	{ID: 4, Name: "明朝精锐", Lv: 100, T1: 20301, TM1: "11378;11379;11380;11381;11382", T2: 20302, TM2: "11378;11379;11380;11381;11382"},
	{ID: 5, Name: "明朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
	{ID: 6, Name: "明朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
	{ID: 7, Name: "宋朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
	{ID: 8, Name: "宋朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
	{ID: 9, Name: "宋朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
	{ID: 10, Name: "宋朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382"},
}

func TestFiletoRows(t *testing.T) {
	_, err := excel.FiletoRows("./testdata/config.xlsx")
	assert.Equal(t, err, nil, err)
	// assert.DeepEqual(t, rows, tmpConfigRows)
}

func TestFiletoSlice(t *testing.T) {
	s, err := excel.FiletoSlice[*tmpExcelConfig]("./testdata/config.xlsx")
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, s, wantExcelConfig)
}

func TestFiletoMap(t *testing.T) {
	m, err := excel.FiletoMap[tmpExcelConfig]("./testdata/config.xlsx", "ID")
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, m, map[string]tmpExcelConfig{
		"1":  {ID: 1, Name: "唐朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		"2":  {ID: 2, Name: "唐朝精锐", Lv: 100, T1: 20301, TM1: "11378;11379;11380;11381;11382", T2: 20302, TM2: "11378;11379;11380;11381;11382"},
		"3":  {ID: 3, Name: "唐朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		"4":  {ID: 4, Name: "明朝精锐", Lv: 100, T1: 20301, TM1: "11378;11379;11380;11381;11382", T2: 20302, TM2: "11378;11379;11380;11381;11382"},
		"5":  {ID: 5, Name: "明朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		"6":  {ID: 6, Name: "明朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		"7":  {ID: 7, Name: "宋朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		"8":  {ID: 8, Name: "宋朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		"9":  {ID: 9, Name: "宋朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		"10": {ID: 10, Name: "宋朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382"},
	})
}

func TestExcelDecoder(t *testing.T) {
	data, err := ioutil.ReadFile("./testdata/config.xlsx")
	assert.Equal(t, err, nil, err)
	// var decoder excel.Decoder
	rows, err := excel.DefaultDecoder.UnmarshalNew(data)
	assert.Equal(t, err, nil, err)

	rows2, err := excel.FiletoRows("./testdata/config.xlsx")
	assert.Equal(t, err, nil, nil)
	assert.DeepEqual(t, rows, rows2)

	// excel.UnmarshaltoMap[tmpExcelConfig]()
}

func TestReadExcelFile(t *testing.T) {
	_, err := excel.Export_readExcelFile("./testdata/config.txt")
	assert.Equal(t, err, excel.ErrUnsupportSuffix)

	_, err = excel.Export_readExcelFile("./testdata/config")
	assert.Equal(t, err, excel.ErrUnsupportSuffix)

	_, err = excel.Export_readExcelFile("./teatdata/config2.xlsx")
	assert.Assert(t, err != excel.ErrUnsupportSuffix)
	assert.Assert(t, err != nil)

	_, err = excel.Export_readExcelFile("./testdata/config.xlsx")
	assert.Equal(t, err, nil, err)
}
