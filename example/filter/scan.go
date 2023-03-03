package filter

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

func ScanJsEftpaths(folder string) StringMap {
	// m := make(map[string]string)

	sm := StringMap{
		Spine:   make(map[string]int),
		Export:  make(map[string]int),
		Special: make(map[string]int),
	}
	readJstringsInFolder(folder, []string{}, sm)
	return sm
}

func ScanJsImagepaths(folder string) map[string]string {
	return nil
}

func ScanCsdImagepaths(folder string, ignoreFolders ...string) map[string]string {
	m := make(map[string]string)
	for i, path := range ignoreFolders {
		path, _ = filepath.Abs(path)
		ignoreFolders[i] = path + string(filepath.Separator)
	}
	WalkGivenExts(folder, func(path string) error {
		idx := slices.IndexFunc(ignoreFolders, func(ignoreDir string) bool {
			return strings.HasPrefix(path, ignoreDir)
		})
		if idx != -1 {
			return nil
		}

		data, _ := os.ReadFile(path)
		tmparr := bytes.Split(data, []byte{'"'}) //strings.Split(string(data), "\"")
		///
		// arr := make([]string, 0, len(tmparr)/2)
		for i, tmp := range tmparr {
			if i%2 == 1 && (bytes.HasSuffix(tmp, []byte(".png")) || bytes.HasSuffix(tmp, []byte(".jpg"))) {
				//
				str := string(tmp)
				m[str] = filepath.Ext(str)
			}
		}
		return nil
	}, ".csd")
	return m
}

// @return map[path]ext
func ScanFilepaths(folder string) map[string]string {
	m := make(map[string]string)
	filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			m[path] = filepath.Ext(path)[1:]
		}
		return nil
	})
	return m
}
