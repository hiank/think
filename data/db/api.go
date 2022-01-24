package db

import "io"

type KvDB interface {
	Get(k string, v interface{}) (found bool, err error)
	Set(k string, v interface{}) error
	Delete(k string) error
	io.Closer
}
