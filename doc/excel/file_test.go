package excel_test

import (
	"reflect"
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

func TestUnmarshalNewMap(t *testing.T) {
	// m, err := excel.UnmarshalNewMap[tmpExcel](tmpRows, "Lv")
	// assert.Equal(t, err, nil, err)
	// assert.DeepEqual(t, m, map[string]tmpExcel{"12": {"11", 12}, "2": {"1", 2}})

	// m2 := make(map[string]*tmpExcel)
	_, err := excel.UnmarshalNewMap[int, *tmpExcel](tmpRows, "关卡ID")
	assert.Equal(t, err, excel.ErrInvalidKeyType)

	m2, err := excel.UnmarshalNewMap[string, *tmpExcel](tmpRows, "关卡ID")
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, m2, map[string]*tmpExcel{"11": {ID: "11", Lv: 12}, "1": {ID: "1", Lv: 2}})
}

func TestUnmarshalNewSlice(t *testing.T) {
	// s, err := excel.UnmarshalNewSlice[*tmpExcel](tmpRows)
	// assert.Equal(t, err, nil, err)
	// assert.DeepEqual(t, s, []*tmpExcel{{"11", 12}, {"1", 2}})

	// var out []tmpExcel
	s, err := excel.UnmarshalNewSlice[tmpExcel](tmpRows)
	assert.Equal(t, err, nil)
	assert.DeepEqual(t, s, []tmpExcel{{ID: "11", Lv: 12}, {ID: "1", Lv: 2}})

	s1, err := excel.UnmarshalNewSlice[*tmpExcel](tmpRows)
	assert.Equal(t, err, nil)
	assert.DeepEqual(t, s1, []*tmpExcel{{ID: "11", Lv: 12}, {ID: "1", Lv: 2}})
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
	ID    uint   `excel:"关卡ID"`
	Name  string `excel:"关卡名"`
	Lv    uint   `excel:"怪物等级"`
	T1    uint   `excel:"队伍1神兽"`
	TM1   string `excel:"队伍1"`
	T2    uint   `excel:"队伍2神兽"`
	TM2   string `excel:"队伍2"`
	NOTAG string
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

func TestReadFileToSlice(t *testing.T) {
	s, err := excel.ReadFileNewSlice[*tmpExcelConfig]("./testdata/config.xlsx")
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, s, wantExcelConfig)
}

func TestReadFileNewMap(t *testing.T) {
	m, err := excel.ReadFileNewMap[uint, tmpExcelConfig]("./testdata/config.xlsx", "Non")
	assert.Equal(t, err, excel.ErrNonKeyField)
	assert.Equal(t, len(m), 0)

	m, err = excel.ReadFileNewMap[uint, tmpExcelConfig]("./testdata/config.xlsx", "关卡ID")
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, m, map[uint]tmpExcelConfig{
		1:  {ID: 1, Name: "唐朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		2:  {ID: 2, Name: "唐朝精锐", Lv: 100, T1: 20301, TM1: "11378;11379;11380;11381;11382", T2: 20302, TM2: "11378;11379;11380;11381;11382"},
		3:  {ID: 3, Name: "唐朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		4:  {ID: 4, Name: "明朝精锐", Lv: 100, T1: 20301, TM1: "11378;11379;11380;11381;11382", T2: 20302, TM2: "11378;11379;11380;11381;11382"},
		5:  {ID: 5, Name: "明朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		6:  {ID: 6, Name: "明朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		7:  {ID: 7, Name: "宋朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		8:  {ID: 8, Name: "宋朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		9:  {ID: 9, Name: "宋朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382", T2: 20304, TM2: "11378;11379;11380;11381;11382"},
		10: {ID: 10, Name: "宋朝精锐", Lv: 100, T1: 20303, TM1: "11378;11379;11380;11381;11382"},
	})
}

// func TestExcelDecoder(t *testing.T) {
// 	data, err := os.ReadFile("./testdata/config.xlsx")
// 	assert.Equal(t, err, nil, err)
// 	// var decoder excel.Decoder
// 	rows, err := excel.DefaultDecoder.UnmarshalNew(data)
// 	assert.Equal(t, err, nil, err)

// 	rows2, err := excel.Export_readFileToRows("./testdata/config.xlsx")
// 	assert.Equal(t, err, nil, nil)
// 	assert.DeepEqual(t, rows, rows2)

// 	// excel.UnmarshaltoMap[tmpExcelConfig]()
// }

func TestReadFileToRows(t *testing.T) {
	_, err := excel.Export_readFileToRows("./testdata/config.txt")
	assert.Equal(t, err, excel.ErrUnsupportFileExt)

	_, err = excel.Export_readFileToRows("./testdata/config")
	assert.Equal(t, err, excel.ErrUnsupportFileExt)

	_, err = excel.Export_readFileToRows("./teatdata/config2.xlsx")
	assert.Assert(t, err != excel.ErrUnsupportFileExt)
	assert.Assert(t, err != nil)

	_, err = excel.Export_readFileToRows("./testdata/config.xlsx")
	assert.Equal(t, err, nil, err)

	_, err = excel.Export_readFileToRows("./testdata/unexist.xlsx")
	assert.Assert(t, err != nil)

	_, err = excel.Export_readFileToRows("./testdata/empty.xlsx")
	assert.Equal(t, err, excel.ErrFailedToReadRows, "empty excel")
}

func TestReflectFiled(t *testing.T) {
	var v tmpExcel
	rt := reflect.TypeOf(&v)
	///
	// ktag := "ID"
	num := rt.Elem().NumField()
	assert.Equal(t, num, 3)

	rv := reflect.ValueOf(v)
	num = rv.NumField()
	assert.Equal(t, num, 3)
}

func TestTitle(t *testing.T) {
	// ID   uint   `excel:"关卡ID"`
	// Name string `excel:"关卡名"`
	// Lv   uint   `excel:"怪物等级"`
	// T1   uint   `excel:"队伍1神兽"`
	// TM1  string `excel:"队伍1"`
	// T2   uint   `excel:"队伍2神兽"`
	// TM2  string `excel:"队伍2"`
	row := []string{"怪物等级", "Non", "关卡ID", "队伍1", "T2", "NOTAG"}
	var intv int = 12
	_, err := excel.Export_newTitle(row, reflect.TypeOf(intv))
	assert.Equal(t, err, excel.ErrNotStruct)
	_, err = excel.Export_newTitle(row, reflect.TypeOf(&intv))
	assert.Equal(t, err, excel.ErrNotStruct)

	var tec tmpExcelConfig
	title, err := excel.Export_newTitle(row, reflect.TypeOf(tec))
	///
	assert.Equal(t, err, nil)
	// assert.Assert(t, title != nil)
	r, m := excel.ExportGetTitleMember(title)
	assert.DeepEqual(t, r, row)
	assert.DeepEqual(t, m, map[int]string{
		0: "Lv",
		2: "ID",
		3: "TM1",
		// 4: "T2",
		5: "NOTAG",
	})
}

func TestRows(t *testing.T) {

	_, err := excel.Export_newRows[*tmpExcel](tmpRows[:1])
	assert.Equal(t, err, excel.ErrInvalidParam)
	//
	r, err := excel.Export_newRows[*tmpExcel](tmpRows)
	assert.Equal(t, err, nil)
	_, found := r.GetFieldName("Name")
	assert.Assert(t, !found)
	fn, found := r.GetFieldName("名字")
	assert.Assert(t, found)
	assert.Equal(t, fn, "Name")

	// err := excel.Export_rangeRows(tmpRows, func(h *excel.Header[tmpExcel], s []string) error {
	// 	return io.EOF
	// })
	// assert.Equal(t, err, io.EOF)

	// cnt := 0
	// err = excel.Export_rangeRows(tmpRows, func(h *excel.Header[*tmpExcel], s []string) error {
	// 	cnt++
	// 	return nil
	// })
	// assert.Equal(t, err, nil, err)
	// assert.Equal(t, cnt, len(tmpRows)-1)

	// cnt = 0
	// err = excel.Export_rangeRows([][]string{tmpRows[0]}, func(h *excel.Header[*tmpExcel], s []string) error {
	// 	cnt++
	// 	return nil
	// })
	// assert.Equal(t, err, excel.ErrInvalidParam)
	// assert.Equal(t, cnt, 0)

	// err = excel.Export_rangeRows([][]string{}, func(h *excel.Header[*tmpExcel], s []string) error {
	// 	cnt++
	// 	return nil
	// })
	// assert.Equal(t, err, excel.ErrInvalidParam)
	// assert.Equal(t, cnt, 0)
}
