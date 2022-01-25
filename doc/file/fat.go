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

// func (f *fat) Encode(v interface{}) error {
// 	return nil
// }

func (f *fat) Decode(outVals ...interface{}) (err error) {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, buffer := range f.m {
		err = pushError(err, buffer.Decode(outVals...))
	}
	return
}

func (f *fat) Val() []byte {
	return []byte("not support")
}

// //ParseAndClear parse loaded data to configs
// //clear loaded data at the end
// func (lm *liteMux) ParseAndClear(configs ...interface{}) {
// 	lm.mux.Lock()
// 	defer lm.mux.Unlock()
// 	for em := lm.list.Front(); em != nil; em = em.Next() {
// 		// em.Value.(parser).parse(configs...)
// 		for _, cfg := range configs {
// 			em.Value.(doc.Decoder).Decode(cfg)
// 		}
// 	}
// 	lm.loaded, lm.list = make(map[string]byte), list.New()
// }
