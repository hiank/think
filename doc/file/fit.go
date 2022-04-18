package file

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/hiank/think/doc"
)

//fit buffer
//only one valid data (conver old loaded data)
//not goroutine safe
type fit struct {
	form Form
	// doc  doc.Doc
	b *doc.B
}

func (f *fit) LoadFile(paths ...string) error {
	done := fmt.Errorf("load done")
	for _, path := range paths {
		if filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				if path, _ = filepath.Abs(path); f.load(path) {
					//load success then skip walk
					err = done
				}
			}
			return err
		}) == done {
			return nil
		}
	}
	return fmt.Errorf("non supporting file")
}

func (f *fit) load(path string) (suc bool) {
	if form := pathToForm(path); form == f.form {
		if data, err := ioutil.ReadFile(path); err == nil {
			suc = f.LoadBytes(form, data) == nil
		}
	}
	return
}

//LoadBytes load given form bytes value
//the value will cover the old bytes value
func (f *fit) LoadBytes(form Form, vals ...[]byte) (err error) {
	if form != f.form {
		return fmt.Errorf("only support Form (%d), but given Form (%d)", f.form, form)
	}
	if len(vals) == 0 {
		return fmt.Errorf("non bytes value passed")
	}
	return f.b.Encode(vals[0])
}

func (f *fit) Decode(outVals ...any) (err error) {
	for _, out := range outVals {
		err = pushError(err, f.b.Decode(out))
	}
	return
}

func (f *fit) Clear() {
	f.b.Encode([]byte{})
}
