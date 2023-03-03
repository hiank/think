package filter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hiank/think/example/filter"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"gotest.tools/v3/assert"
)

const (
	// jsfolder = "D:/hiank/ws/go/read-csd/src/ducheng/searchSecret/"
	jsfolder = "D:/hiank/svn/tmp/src/"
	// jsfolder = "D:/hiank/svn/tmp/src/view/nationConstruction/cityBattle/"
	eftFolder = "D:/hiank/svn/tmp/res/image/"
	csvFolder = "D:/hiank/svn/tmp/res/config/"
)

var jsIgnores = []string{
	"D:/hiank/svn/tmp/src/resource.js",
	"D:/hiank/svn/tmp/src/utils/jszip.min.js",
	"D:/hiank/svn/tmp/src/common/protobuf.min.js",
}

func TestFilterProtoFromJs(t *testing.T) {
	var jsfolder = "D:/hiank/ws/go/read-csd/src/holiday"

	dst := "./protonames.xml"
	m := make(map[string]int)
	filter.UnmarshalProtoFolder(jsfolder, m)

	arr := maps.Keys(m)
	slices.Sort(arr)

	data := make([]byte, 0, 1024)
	for _, proto := range arr {
		data = append(data, proto...)
		data = append(data, '\n')
	}
	os.WriteFile(dst, data, 0777)
}

func TestFilterFuncidFromJs(t *testing.T) {
	dst := "./funcids.xml"
	m := make(map[string]int)
	filter.UnmarshalFuncidFolder(jsfolder, m)

	arr := maps.Keys(m)
	slices.Sort(arr)

	data := make([]byte, 0, 1024)
	for _, funcid := range arr {
		data = append(data, funcid...)
		data = append(data, '\n')
	}
	os.WriteFile(dst, data, 0777)
}

func TestSlicesContains(t *testing.T) {
	vals := []string{
		"test1",
		"test2",
	}
	suc := slices.Contains(vals, "test1")
	assert.Assert(t, suc)
}

func TestFilterJsEfts(t *testing.T) {
	jsfolder := "D:/hiank/ws/go/read-csd/src/holiday"
	sm := filter.ScanJsEftpaths(jsfolder)

	arr := []map[string]int{
		// sm.Special,
		sm.Spine,
		sm.Export,
	}
	data := make([]byte, 0, 1024)
	for _, m := range arr {
		for name := range m {
			data = append(data, name...)
			data = append(data, '\n')
		}
	}
	os.WriteFile("./efts/jsefts.xml", data, 0777)
}

func TestFilterEfts(t *testing.T) {

	// jsfolder = "D:/hiank/ws/go/read-csd/src/ducheng/searchSecret/"
	jsfolder := "D:/hiank/svn/tmp/src/"
	// jsfolder = "D:/hiank/svn/tmp/src/view/nationConstruction/cityBattle/"
	eftFolder := "D:/hiank/svn/tmp/res/image/"
	csvFolder := "D:/hiank/svn/tmp/res/config/"

	var jsIgnores = []string{
		"D:/hiank/svn/tmp/src/resource.js",
		"D:/hiank/svn/tmp/src/utils/jszip.min.js",
		"D:/hiank/svn/tmp/src/common/protobuf.min.js",
		"D:/hiank/svn/tmp/src/utils/EffectTable.js",
	}
	var eftFolderIgnores = []string{
		"D:/hiank/svn/tmp/res/image/action/wujiang/",
		"D:/hiank/svn/tmp/res/image/action/hongyan/",
	}

	unusedst, notfoundst := "./efts/unused.xml", "./efts/notfound.xml"
	unusedm, notfounds := filter.UnusedEfts(eftFolder, jsfolder, csvFolder, jsIgnores, eftFolderIgnores)
	vals := maps.Values(unusedm)
	slices.Sort(vals)
	data := make([]byte, 0, 1024)
	for _, name := range vals {
		data = append(data, name...)
		data = append(data, '\n')
	}
	os.WriteFile(unusedst, data, 0777)

	data = make([]byte, 0, 1024)
	slices.Sort(notfounds)
	for _, name := range notfounds {
		data = append(data, name...)
		data = append(data, '\n')
	}
	os.WriteFile(notfoundst, data, 0777)
}

func TestScanCsvJsAndFiles(t *testing.T) {
	jsPath := "D:/hiank/svn/tmp/src/data/Config.js"
	csvFolder := "D:/hiank/svn/gj_mori4_develop/doc/数据表/1.0.0版本/csv"
	csvFileM := make(map[string]string)
	filter.ScanCsvInFolder(csvFolder, csvFileM)
	filenames, texts := filter.ScanCsvJsAndFiles(jsPath, "./config/csv", csvFileM)

	data := make([]byte, 0, 1024)
	// slices.Sort(notfounds)
	for _, name := range filenames {
		data = append(data, name...)
		data = append(data, '\n')
	}
	os.WriteFile("./config/unused.xml", data, 0777)

	data = make([]byte, 0, 1024)
	// slices.Sort(notfounds)
	for _, name := range texts {
		data = append(data, name...)
		data = append(data, '\n')
	}
	os.WriteFile("./config/nofound.xml", data, 0777)

	jsFolder := "D:/hiank/svn/tmp/src"
	jsIgnores := []string{
		"D:/hiank/svn/tmp/src/data/Config.js",
		"D:/hiank/svn/tmp/src/utils/jszip.min.js",
		"D:/hiank/svn/tmp/src/common/protobuf.min.js",
	}

	notused := filter.ScanNotUsedInJs(jsFolder, csvFileM, jsIgnores)
	keys := maps.Keys(notused)
	slices.Sort(keys)

	data = make([]byte, 0, 1024)
	// slices.Sort(notfounds)
	for _, name := range keys {
		data = append(data, name...)
		data = append(data, '\n')
	}
	os.WriteFile("./config/nofoundinjs.xml", data, 0777)

}

func TestUnmarshalHudFolder(t *testing.T) {
	jsFolder := "D:/hiank/svn/gj_mori4_develop/client/branches/gj_mori4_1.0.0/src"
	csdFolder := "D:/hiank/svn/gj_mori4_develop/client/branches/gj_mori4_1.0.0/cocosstudio/ccs"
	m, mcsd := make(map[string]int), make(map[string]string)
	filter.ListHUD_LISTUsed(jsFolder, m)

	filter.WalkGivenExts(csdFolder, func(path string) error {
		///
		name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		if _, ok := m[name]; !ok {
			mcsd[name] = path
		} else {
			delete(m, name)
		}
		return nil
	}, ".csd")

	t.Logf("%v", m)

	keys := maps.Keys(m)
	slices.Sort(keys)

	data := make([]byte, 0, 1024)
	// slices.Sort(notfounds)
	for _, name := range keys {
		data = append(data, name...)
		data = append(data, '\n')
	}
	os.WriteFile("./nocsd.txt", data, 0777)

	values := maps.Values(mcsd)
	slices.Sort(values)

	data = make([]byte, 0, 1024)
	// slices.Sort(notfounds)
	for _, path := range values {
		data = append(data, path...)
		data = append(data, '\n')
	}
	os.WriteFile("./notusedcsd.txt", data, 0777)

}

func TestFilterCsdImage(t *testing.T) {
	csdFolder := "D:/hiank/svn/gj_mori4_develop/client/branches/gj_mori4_1.0.0/cocosstudio/ccs"
	uiRoot := "D:/hiank/svn/gj_mori4_develop/art/1.0.0版本/"
	uiFolder := uiRoot + "res/image/ui"

	// t.Log(csdFolder)
	mui := filter.ScanFilepaths(uiFolder)
	mcsd := filter.ScanCsdImagepaths(csdFolder)

	// t.Log(mcsd)

	///
	uiRoot, _ = filepath.Abs(uiRoot)
	uiRoot += string(filepath.Separator)
	rootLen := len(uiRoot)

	t.Log(len(mui), len(mcsd))
	///
	for path := range mui {
		if strings.HasPrefix(path, uiRoot) {
			key := strings.ReplaceAll(path[rootLen:], string(filepath.Separator), "/")
			if _, ok := mcsd[key]; ok {
				delete(mui, path)
				delete(mcsd, key)
			}
		}
	}
	t.Log("....", len(mui), len(mcsd))

	paths := maps.Keys(mcsd)
	slices.Sort(paths)

	data := make([]byte, 0, 1024)
	for _, path := range paths {
		data = append(data, path...)
		data = append(data, '\n')
	}
	os.WriteFile("./tmp/csd-nonart.txt", data, 0777)

	paths = maps.Keys(mui)
	slices.Sort(paths)

	data = make([]byte, 0, 1024)
	for _, path := range paths {
		data = append(data, path...)
		data = append(data, '\n')
	}
	os.WriteFile("./tmp/art-noncsd.txt", data, 0777)
}

func TestScanFilepaths(t *testing.T) {
	folder := "D:/hiank/svn/gj_mori4_develop/art/1.0.0版本/res/image"
	m := filter.ScanFilepaths(folder)
	///
	mm := make(map[string]map[string]int)
	for path, ext := range m {
		tm, ok := mm[ext]
		if !ok {
			tm = make(map[string]int)
			mm[ext] = tm
		}
		tm[path] = 1
	}
	///
	var root = "./tmp/"
	for ext, m := range mm {
		paths := maps.Keys(m)
		slices.Sort(paths)

		data := make([]byte, 0, 1024)
		for _, path := range paths {
			data = append(data, path...)
			data = append(data, '\n')
		}
		os.WriteFile(root+ext+".txt", data, 0777)
	}
}

func TestScanCsdImagepaths(t *testing.T) {
	//
	folder, root := "D:/hiank/svn/gj_mori4_develop/client/branches/gj_mori4_1.0.0/cocosstudio/ccs", "./tmp/"
	m := filter.ScanCsdImagepaths(folder)
	paths := maps.Keys(m)
	slices.Sort(paths)

	data := make([]byte, 0, 1024)
	for _, path := range paths {
		data = append(data, path...)
		data = append(data, '\n')
	}
	os.WriteFile(root+"csdimages.txt", data, 0777)
}

func TestFilepath(t *testing.T) {
	path := "tmp/test/"
	//
	path, _ = filepath.Abs(path)
	t.Log(path)
}


// func TestFilter