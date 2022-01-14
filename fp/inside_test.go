package fp

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"gotest.tools/v3/assert"
)

type testConfig struct {
	Tik string `json:"sys.Tik"`
	Tok int    `yaml:"Tok"`
}

func TestJsonParser(t *testing.T) {
	u := &jsonData{data: []byte(`{"sys.Tik": "nil"}`)}
	var cfg testConfig
	u.parse(&cfg)

	assert.Equal(t, cfg.Tik, "nil")

	u = &jsonData{data: []byte(`{"sys.Tik": "overwrite"}`)}
	u.parse(&cfg)
	assert.Equal(t, cfg.Tik, "overwrite", "overwrite previous settings")
}

func TestYamlParser(t *testing.T) {
	u := &yamlData{data: []byte(`Tok: 2`)}
	var cfg testConfig
	u.parse(&cfg)

	assert.Equal(t, cfg.Tok, 2)

	u = &yamlData{data: []byte(`Tok: 3`)}
	u.parse(&cfg)
	assert.Equal(t, cfg.Tok, 3, "overwrite previous settings")
}

func TestMarch(t *testing.T) {
	paths := match("testdata")
	assert.Equal(t, len(paths), 5)

	root, _ := filepath.Abs("testdata")
	// strings.
	sp := string(filepath.Separator)
	root += sp
	// t.Log(root)
	names := []string{
		"config.json",
		"config.yaml",
		"dep" + sp + "dep.YaMl",
		"dep" + sp + "dep.json", //NOTE: 'j' > 'Y'
		"dep2" + sp + "dep.json",
	}
	for i, path := range paths {
		assert.Equal(t, path, root+names[i])
	}

	path := match("testdata/config.json")[0]
	assert.Equal(t, path, paths[0])

	paths = match("testdata/non.json")
	assert.Equal(t, len(paths), 0)

	_, err := ioutil.ReadDir("testdata")
	assert.Assert(t, err == nil)
}

// func TestReadExcel(t *testing.T) {
// 	// u := fp.NewParser()
// 	// u.LoadFile("testdata/config.xls")
// 	f, err := excelize.OpenFile("testdata/config.xlsx")
// 	assert.Assert(t, err == nil)
// 	defer f.Close()
// 	// f.GetSheet
// 	val, _ := f.GetRows(f.GetSheetList()[0])
// 	t.Log(val)
// }

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

func TestReflectType(t *testing.T) {
	val := &testExcel{}
	fv := reflect.ValueOf(val)
	assert.Equal(t, fv.Kind(), reflect.Ptr)

	fv = fv.Elem()
	assert.Equal(t, fv.Kind(), reflect.Struct)

	ff, _ := fv.Type().FieldByName("ID")
	tag := ff.Tag.Get("excel")
	assert.Equal(t, tag, "关卡ID")

	assert.Equal(t, ff.Name, "ID")

	rc := newRowsConv(testRows, reflect.TypeOf(*val))
	vals := rc.Unmarshal() //unmarshalRows(testRows, reflect.TypeOf(*val))
	assert.Equal(t, len(vals), 2)
	// assert.Equal(t, vals[])
	val = vals[0].(*testExcel)
	assert.Equal(t, val.ID, "11")
}
