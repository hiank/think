package filter

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestUnmarshalToStringMap(t *testing.T) {
	sm := StringMap{
		Export:  map[string]int{},
		Spine:   map[string]int{},
		Special: map[string]int{},
	}
	unmarshalToStringMap(sm, "res/image/effects/eft_gongcheng_zhentexiao/zhentexiao1.ExportJson")
	assert.Equal(t, len(sm.Special), 4)
	assert.Equal(t, len(sm.Spine), 0)
	assert.Equal(t, len(sm.Export), 1)
	// t.Logf("%v", sm.Special)
	assert.Equal(t, sm.Special["zhentexiao1"], 1)
}

func TestReadJs(t *testing.T) {
	// path := "D:/hiank/svn/tmp/src/view/nationConstruction/cityBattle/CityBattleLayer.js"
	// path := "D:/hiank/svn/tmp/src/extern/UIExtern.js"
	// path := "D:/hiank/svn/tmp/src/manager/mSDK.js"
	// path := "D:/hiank/svn/tmp/src/common/protobuf.min.js"
	paths := []string{
		// "D:/hiank/svn/tmp/src/view/nationConstruction/cityBattle/CityBattleLayer.js",
		// "D:/hiank/svn/tmp/src/extern/UIExtern.js",
		// "D:/hiank/svn/tmp/src/manager/mSDK.js",
		// "D:/hiank/svn/tmp/src/common/protobuf.min.js",
		"D:/hiank/svn/tmp/src/utils/jszip.min.js",
	}

	for _, path := range paths {
		bs := ReadJs(path)
		// t.Log(string(bs))

		cnt := bytes.Count(bs, []byte{'"'})
		t.Log("+++++++", cnt)

		// os.WriteFile("./tmp.js", bs, 0777)
	}

}

func TestReadJstringsInFile(t *testing.T) {
	// path := "D:/hiank/svn/tmp/src/extern/UIExtern.js"
	path := "D:/hiank/svn/tmp/src/common/protobuf.min.js"
	sm := StringMap{
		Export:  map[string]int{},
		Spine:   map[string]int{},
		Special: map[string]int{},
	}
	readJstringsInFile(path, sm)
	t.Log(sm)
}

func TestFilepathDir(t *testing.T) {
	path := "D:/hiank/svn/tmp/src/common/protobuf.min.js"
	dir := filepath.Dir(path)
	absDir, _ := filepath.Abs("D:/hiank/svn/tmp/src/common")
	assert.Equal(t, dir, absDir)

	// filepath.
	base := filepath.Base(path)
	t.Log(base)
	ext := filepath.Ext(path)
	t.Log(ext)
	blob, _ := filepath.Glob(path)
	t.Log(blob)
	absDir, _ = filepath.Abs("D:/hiank/svn/tmp")
	// hp := filepath.HasPrefix(path, absDir)
	// t.Log(hp)
	slash := filepath.FromSlash(path)
	t.Log(slash)
	matched, _ := filepath.Match(absDir, path)
	t.Log(matched)

	absPath, _ := filepath.Abs(path)
	// strings.HasPrefix(absPath, absDir)
	dirLen := len(absDir)
	ru, _, _ := strings.NewReader(absPath[dirLen:]).ReadRune()
	assert.Equal(t, ru, filepath.Separator)
	// absPath[dirLen]
	// assert.Equal(t, abs)
	assert.Equal(t, absPath[:dirLen], absDir)

	path1, _ := filepath.Abs("D:/hiank/")
	path2, _ := filepath.Abs("D:/hiank")
	assert.Equal(t, path1, path2)
	assert.Equal(t, path1[len(path1)-1], byte('k'), "not include separator")
	// filepath.Rel()
	// filepath.Split()
	// filepath.

	st, _ := os.Stat("./testdata/empty")
	t.Log(st.Size())

	st, _ = os.Stat("./testdata/folder")
	t.Log(st.Size())
	// st.Size()

	fis, _ := ioutil.ReadDir("./testdata/folder")
	assert.Equal(t, len(fis), 1)

	fis, _ = ioutil.ReadDir("./testdata/empty")
	assert.Equal(t, len(fis), 0)
}

func TestHasPrefixDir(t *testing.T) {
	has := hasPrefixDir("D:/hiank/", "D:/hiank/tmp/src/txt.js")
	assert.Assert(t, has)

	has = hasPrefixDir("D:/hi", "D:/hiank/tmp/src/txt.js")
	assert.Assert(t, !has)
}

func TestScanCsvInJs(t *testing.T) {
	m := map[string]string{
		"XZmYlcZhuanPanRewared": "XZmYlcZhuanPanRewared.js",
	}
	// out := make(map[string]string)
	scanCsvInJs("D:/hiank/svn/tmp/src/view/entertainmentCity/PrizeWheel.js", m)
	assert.Equal(t, len(m), 0)
	// assert.Equal(t, len(out), 1)
}

// func TestSplitRightMark(t *testing.T) {
// 	arr := splitRightMark([]byte(`"test`))
// 	assert.Equal(t, len(arr[0]), 0)
// 	assert.DeepEqual(t, arr[1], []byte("test"))

// 	arr = splitRightMark([]byte(`left"right"end"`))
// 	assert.DeepEqual(t, arr[0], []byte("left"))
// 	assert.DeepEqual(t, arr[1], []byte(`right"end"`))
// }
