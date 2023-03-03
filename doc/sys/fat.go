package sys

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/hiank/think/doc"
	"golang.org/x/exp/slices"
	"k8s.io/klog/v2"
)

type coderHub struct {
	loaded []string //[]path
	coders []doc.Coder
}

func (ch *coderHub) loadFile(path string) {
	if !slices.Contains(ch.loaded, path) {
		//
		coder, err := doc.ReadFile(path)
		if err != nil {
			klog.Warning("settings:", err)
			return
		}
		ch.loaded = append(ch.loaded, path)
		ch.coders = append(ch.coders, coder)
	}
}

// Fat fat parser. unmarshal []byte to passed value
type Fat struct {
	mux sync.RWMutex
	hub *coderHub
}

func NewFat() *Fat {
	return &Fat{
		hub: &coderHub{
			loaded: make([]string, 0),
			coders: make([]doc.Coder, 0),
		},
	}
}

// Load load files to coder
// folder or file
func (fat *Fat) Load(paths ...string) {
	fat.mux.Lock()
	defer fat.mux.Unlock()
	for _, path := range paths {
		filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				fat.hub.loadFile(path)
			}
			return nil
		})
	}
}

func (fat *Fat) LoadBytes(data []byte, f doc.Format) error {
	///
	var coder doc.Coder
	switch f {
	case doc.FormatJson:
		var jd doc.Json = data
		coder = &jd
	case doc.FormatYaml:
		var yd doc.Yaml = data
		coder = &yd
	default:
		return doc.ErrUnsupportFormat
	}
	fat.mux.Lock()
	defer fat.mux.Unlock()
	////
	fat.hub.coders = append(fat.hub.coders, coder)
	return nil
}

// UnmarshalTo unmarshal cached *Bytes to vals
// it will stop work when first unmarshal fialed and return the unmarshal error
func (fat *Fat) UnmarshalTo(vals ...any) (err error) {
	fat.mux.RLock()
	defer fat.mux.RUnlock()
	for _, coder := range fat.hub.coders {
		for _, val := range vals {
			if derr := coder.Decode(val); derr != nil {
				klog.Warning(derr)
				err = fmt.Errorf("%v: %v", err, derr)
			}
		}
	}
	return
}

func (fat *Fat) Release() {
	*fat = *NewFat()
}
