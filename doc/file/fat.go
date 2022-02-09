package file

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sync"

	"k8s.io/klog/v2"
)

const (
	bkRule = `_bytes_key_%d`
)

type fat struct {
	mux sync.Mutex
	num int                //bytes loaded count
	m   map[string]Decoder //sync.Map //map[path/byteskey]Buffer
}

//LoadFile load folder|.json|.yaml
//read all file contents sync to lm.list
func (f *fat) LoadFile(paths ...string) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, path := range paths {
		filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				path, _ = filepath.Abs(path)
				f.load(path)
			}
			return err
		})
	}
	return nil
}

func (f *fat) load(path string) {
	if _, ok := f.m[path]; ok {
		return ///stored
	}
	form := pathToForm(path)
	d := Fit(form)
	if d == nil {
		klog.Warning("not support file", path)
		return ///not support the form
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		klog.Warning(err)
		return
	}
	d.LoadBytes(form, buf)
	f.m[path] = d
}

func (f *fat) LoadBytes(form Form, vals ...[]byte) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, v := range vals {
		buffer := Fit(form)
		if buffer == nil {
			return fmt.Errorf("invalid Form: not support form (%d)", form)
		}
		buffer.LoadBytes(form, v)
		f.m[fmt.Sprintf(bkRule, f.num)] = buffer
		f.num++
	}
	return nil
}

func (f *fat) Decode(outVals ...interface{}) (err error) {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, buffer := range f.m {
		err = pushError(err, buffer.Decode(outVals...))
	}
	return
}

func (f *fat) Clear() {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, decoder := range f.m {
		decoder.Clear()
	}
	f.m, f.num = make(map[string]Decoder), 0
}
