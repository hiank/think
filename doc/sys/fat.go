package sys

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sync"

	"k8s.io/klog/v2"
)

//Fat fat parser. unmarshal []byte to passed value
type Fat struct {
	mux    sync.Mutex
	loaded map[string]bool
	pool   []*Bytes
}

func NewFat() *Fat {
	return &Fat{
		loaded: make(map[string]bool),
		pool:   make([]*Bytes, 0, 8),
	}
}

//LoadFile
func (fat *Fat) LoadFiles(paths ...string) {
	fat.mux.Lock()
	defer fat.mux.Unlock()
	for _, path := range paths {
		filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				path, _ = filepath.Abs(path)
				fat.loadFile(path)
			}
			return err
		})
	}
}

func (fat *Fat) loadFile(path string) {
	if _, ok := fat.loaded[path]; ok {
		return ///stored
	}
	fat.loaded[path] = true
	err := fat.loadBytes(formatFromPath(path), func() ([]byte, error) { return ioutil.ReadFile(path) })
	if err != nil {
		klog.Warning(err)
	}
}

func (fat *Fat) loadBytes(f Format, get func() ([]byte, error)) error {
	b, err := formatoBytes(f, get)
	if err == nil {
		fat.pool = append(fat.pool, b)
	}
	return err
}

func (fat *Fat) LoadBytes(data []byte, f Format) error {
	fat.mux.Lock()
	defer fat.mux.Unlock()
	return fat.loadBytes(f, func() ([]byte, error) { return data, nil })
}

//UnmarshalTo unmarshal cached *Bytes to vals
//it will stop work when first unmarshal fialed and return the unmarshal error
func (fat *Fat) UnmarshalTo(vals ...any) (err error) {
L:
	for _, b := range fat.pool {
		for _, val := range vals {
			if err = b.UnmarshalTo(val); err != nil {
				break L
			}
		}
	}
	return
}

func (fat *Fat) Release() {
	*fat = *NewFat()
}
