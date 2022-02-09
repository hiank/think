package run

import (
	"errors"
	"io"
	"sync"
)

type onceCloser struct {
	once sync.Once
	f    func() error
}

func NewOnceCloser(f func() error) io.Closer {
	return &onceCloser{f: f}
}

func (oc *onceCloser) Close() (err error) {
	err = errors.New("closed")
	oc.once.Do(func() {
		err = oc.f()
	})
	return
}
