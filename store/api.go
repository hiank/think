package store

import "io"

type Dictionary[KT ~string] interface {
	Scan(k KT, out any) (found bool, err error)
	Set(k KT, v any) error
	Del(k KT, out ...any) error
	io.Closer
}

type EasyDictionary Dictionary[string]
