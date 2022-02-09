package net

import (
	"reflect"
	"sync"

	"k8s.io/klog/v2"
)

const (
	DefaultHandler string = ""
)

type fathandler struct {
	m sync.Map
	// kd KeyDecoder
}

//Handle add handler for message recv
//use k's Type Name as key
func (fh *fathandler) Handle(k interface{}, h Handler) {
	sk, ok := k.(string)
	if !ok {
		rv := reflect.ValueOf(k)
		for rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		sk = rv.Type().Name()
	}
	fh.m.Store(sk, h)
}

//Process message
func (fh *fathandler) Process(id string, d *Doc) {
	mv, loaded := fh.m.Load(d.TypeName())
	if !loaded {
		if mv, loaded = fh.m.Load(DefaultHandler); !loaded {
			klog.Warning("cannot find handler for handle message recv by conn")
			return
		}
	}
	mv.(Handler).Process(id, d)
}
