package set

import (
	"container/list"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/klog/v2"
)

//LookRootFolders 检查并返回所有'根'目录
//如果路径不是目录，舍弃
//如果路径包含于另一个路径下，舍弃
func LookRootFolders(folders []string) []string {
	bak := make([]string, 0, len(folders))
L:
	for _, folder := range folders {
		folder, err := filepath.Abs(folder)
		if err != nil {
			continue L
		}
		if info, err := os.Stat(folder); err != nil || info.IsDir() {
			continue L
		}
		for i, path := range bak {
			switch {
			case strings.Index(folder, path+string(os.PathSeparator)) == 0:
				continue L
			case strings.Index(path, folder+string(os.PathSeparator)) == 0:
				bak[i] = folder
				continue L
			}
		}
		bak = append(bak, folder)
	}
	return bak
}

//WalkText 读取路径下所有指定后缀的文件内容到cacheMap中
func WalkText(cacheMap map[string]*list.List, dir string, suffixes ...string) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		klog.Infof("walk folder error: %v\n", err)
		return
	}
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if suffix, ok := withSuffixes(info.Name(), suffixes); ok {
			cacheText(cacheMap, suffix, path)
		}
		return nil
	})
}

//withSuffixes 判断名称中是否包含指定的一组结尾之一，如果包含，则返回之
func withSuffixes(name string, suffixes []string) (string, bool) {
	idx := strings.LastIndexByte(name, '.')
	if idx == -1 {
		return "", false
	}
	suffix := name[idx+1:]
	for _, wantSuffix := range suffixes {
		if wantSuffix == suffix {
			return suffix, true
		}
	}
	return "", false
}

func cacheText(cacheMap map[string]*list.List, key string, path string) {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		klog.Warning(path, err)
		return
	}
	cache, ok := cacheMap[key]
	if !ok {
		cache = list.New()
		cacheMap[key] = cache
	}
	cache.PushBack(text)
}
