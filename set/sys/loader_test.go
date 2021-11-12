package sys_test

import (
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/hiank/think/set/sys"
	"gotest.tools/v3/assert"
	"k8s.io/klog/v2"
)

var (
	testRootDir = "d:\\hiank\\env"
)

func BenchmarkRecursionTailcall(t *testing.B) {
	loader := sys.Export_newLoader(testRootDir, sys.LTFolder)
	// t.Log(len(arr.Match()))
	loader.Match()
	// t.Log(len(fpaths))
}

func BenchmarkRecursion(t *testing.B) {
	loader := sys.Export_newLoader(testRootDir, sys.LTFolder)
	var match func(string) []string
	match = func(dpath string) []string {
		dpaths, fpaths := loader.ListPaths(dpath)
		for _, dpath := range dpaths {
			fpaths = append(fpaths, match(dpath)...)
		}
		return fpaths
	}
	match(testRootDir)
}

func BenchmarkFilepathWalk(t *testing.B) {
	dpaths, fpaths := make([]string, 0), make([]string, 0)
	filepath.Walk(testRootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			klog.Warningf("sys: Loader-listPaths: Walk for file %s %v", err)
			return err
		}
		if info.IsDir() {
			dpaths = append(dpaths, path)
		} else {
			fpaths = append(fpaths, path)
		}
		return nil
	})
}

type testConf struct {
	Limit int    `json:"sys.Limit"`
	Key   string `yaml:"a"`
}

func TestLoad(t *testing.T) {
	t.Run("Json File", func(t *testing.T) {
		tc := &testConf{}
		loader := sys.Export_newLoader("testdata/config.json", sys.LTFileJson)
		loader.Handle(tc)
		loader.Load()
		assert.Equal(t, tc.Limit, 11)
	})
	t.Run("Yaml File", func(t *testing.T) {
		tc := &testConf{}
		loader := sys.Export_newLoader("testdata/config.yaml", sys.LTFileYaml)
		loader.Handle(tc)
		loader.Load()
		assert.Equal(t, tc.Key, "ws")
	})
	t.Run("Floder", func(t *testing.T) {
		tc := &testConf{}
		loader := sys.Export_newLoader("testdata", sys.LTFolder)
		loader.Handle(tc)
		loader.Load()
		assert.Equal(t, tc.Key, "love-ws")
		assert.Equal(t, tc.Limit, 201)
	})
}

// func TestRecursion(t *testing.T) {
// 	loader := sys.Export_newLoader("d:\\hiank\\env\\protoc", sys.LTFolderDep)
// 	t.Log(len(loader.Match()))
// }

// func TestListPaths(t *testing.T) {
// 	loader := sys.Export_newLoader("d:\\hiank\\env\\protoc", sys.LTFolderDep)
// 	dpaths, fpaths := loader.ListPaths("d:\\hiank\\env\\protoc")
// 	t.Logf("dpaths: %v", dpaths)
// 	t.Logf("fpaths: %v", fpaths)

// 	// filepath.
// }
