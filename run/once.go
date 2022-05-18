package run

import (
	"io"
	"sync"
)

const (
	ErrBeenClosed = Err("has been closed")
)

type onceCloser struct {
	once *sync.Once
	f    func() error
}

func NewOnceCloser(f func() error) io.Closer {
	return &onceCloser{f: f, once: new(sync.Once)}
}

func (oc *onceCloser) Close() (err error) {
	err = ErrBeenClosed
	oc.once.Do(func() {
		err = oc.f()
	})
	return
}
