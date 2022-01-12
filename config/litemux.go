package config

import (
	"container/list"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sync"

	"k8s.io/klog/v2"
)

type liteMux struct {
	list   *list.List
	loaded map[string]byte
	mux    sync.Mutex
}

//NewParser new a IParser
func NewParser() IParser {
	return &liteMux{
		loaded: make(map[string]byte),
		list:   list.New(),
	}
}

//LoadFile load folder|.json|.yaml
//read all file contents sync to lm.list
func (lm *liteMux) LoadFile(paths ...string) {
	lm.mux.Lock()
	defer lm.mux.Unlock()
	for _, path := range paths {
		filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				path, _ = filepath.Abs(path)
				lm.load(path)
			}
			return err
		})
	}
}

func (lm *liteMux) load(path string) {
	if _, ok := lm.loaded[path]; ok {
		return
	}
	lm.loaded[path] = 1
	data, err := ioutil.ReadFile(path)
	if err != nil {
		klog.Warning(err)
		return
	}
	switch suffix(path) {
	case jsonSuffix:
		lm.list.PushBack(&jsonData{data: data})
	case yamlSuffix:
		lm.list.PushBack(&yamlData{data: data})
	default:
		klog.Warningf("not support the file type now: %v", path)
	}
}

//LoadJsonBytes load json data to cache
func (lm *liteMux) LoadJsonBytes(data []byte) {
	lm.mux.Lock()
	defer lm.mux.Unlock()
	lm.list.PushBack(&jsonData{data: data})
}

//LoadYamlBytes load yaml data to cache
func (lm *liteMux) LoadYamlBytes(data []byte) {
	lm.mux.Lock()
	defer lm.mux.Unlock()
	lm.list.PushBack(&yamlData{data: data})
}

//ParseAndClear parse loaded data to configs
//clear loaded data at the end
func (lm *liteMux) ParseAndClear(configs ...IConfig) {
	lm.mux.Lock()
	defer lm.mux.Unlock()
	for em := lm.list.Front(); em != nil; em = em.Next() {
		em.Value.(parser).parse(configs...)
	}
	lm.loaded, lm.list = make(map[string]byte), list.New()
}
