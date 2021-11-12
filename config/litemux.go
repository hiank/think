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

//NewUnmarshaler new a IUnmarshaler
func NewUnmarshaler() IUnmarshaler {
	return &liteMux{
		loaded: make(map[string]byte),
		list:   list.New(),
	}
}

//HandleFile handle folder|.json|.yaml
//read all file contents sync to lm.list
func (lm *liteMux) HandleFile(paths ...string) {
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

//HandleJsonBytes handle json type data
func (lm *liteMux) HandleJsonBytes(data []byte) {
	lm.mux.Lock()
	defer lm.mux.Unlock()
	lm.list.PushBack(&jsonData{data: data})
}

//HandleYamlBytes handle yaml type data
func (lm *liteMux) HandleYamlBytes(data []byte) {
	lm.mux.Lock()
	defer lm.mux.Unlock()
	lm.list.PushBack(&yamlData{data: data})
}

//UnmarshalAndClean unmarshal the handled data to configs
//clean loaded data at the end
func (lm *liteMux) UnmarshalAndClean(configs ...IConfig) {
	lm.mux.Lock()
	defer lm.mux.Unlock()
	for em := lm.list.Front(); em != nil; em = em.Next() {
		em.Value.(unmarshaler).unmarshal(configs...)
	}
	lm.loaded, lm.list = make(map[string]byte), list.New()
}
