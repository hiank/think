package db

import "io"

//KvDB key-value stype database client
type KvDB interface {
	Get(k string, v interface{}) (found bool, err error)
	Set(k string, v interface{}) error
	Delete(k string) error
	io.Closer
}
