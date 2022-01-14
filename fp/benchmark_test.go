package fp

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/klog/v2"
)

var (
	testRootDir = "d:\\hiank\\env"
)

//match 匹配所有文件
//Deprecated: use filepath.WalkDir
func match(path string) (out []string) {
	stat, err := os.Lstat(path)
	switch {
	case err != nil:
		klog.Warning(err)
	case stat.IsDir():
		out = tailMatch(listPaths(path))
	default:
		path, _ = filepath.Abs(path) //NOTE: it will succeed after 'os.Stat' succeed
		out = []string{path}
	}
	return
}

//tailMatch tial call algorithm for match all files
func tailMatch(dpaths, fpaths []string) []string {
	ds, fs := listPaths(dpaths[0])
	dpaths = append(dpaths[1:], ds...)
	fpaths = append(fpaths, fs...)
	if len(dpaths) == 0 {
		return fpaths
	}
	return tailMatch(dpaths, fpaths)
}

//listPaths list folders and files in given folder path
//dpaths director paths
//fpaths file paths
func listPaths(dpath string) (dpaths, fpaths []string) {
	fis, err := ioutil.ReadDir(dpath)
	if err != nil {
		klog.Warning("listPaths: ReadDir ", err)
		return
	}
	fpaths, dpaths = make([]string, 0, len(fis)), make([]string, 0, len(fis))
	dpath, _ = filepath.Abs(dpath) //NOTE: it will succeed after 'ioutil.ReadDir' succeed
	for _, fi := range fis {
		if fi.IsDir() {
			dpaths = append(dpaths, filepath.Join(dpath, fi.Name()))
		} else {
			fpaths = append(fpaths, filepath.Join(dpath, fi.Name()))
		}
	}
	return
}

func BenchmarkRecursionTailcall(t *testing.B) {
	paths := match(testRootDir)
	t.Log(len(paths))
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
	t.Log(len(fpaths), len(dpaths))
}

func BenchmarkFilepathWalkDir(t *testing.B) {
	dpaths, fpaths := make([]string, 0), make([]string, 0)
	filepath.WalkDir(testRootDir, func(path string, info fs.DirEntry, err error) error {
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
	t.Log(len(fpaths), len(dpaths))
}
