package filter

import (
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

func WalkGivenExts(folder string, opt func(path string) error, exts ...string) {
	filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && slices.Contains(exts, filepath.Ext(path)) && !strings.HasSuffix(path, ".min.js") {	//压缩的js文件忽略
			return opt(path)
		}
		return nil
	})
}

// func Scan
