package data

import (
	"github.com/hiank/think/data/db"
)

var (
	Export_KeyTag = func(kt KeyTag) *expKeyTag {
		return &expKeyTag{kt}
	}
	Export_decode      = decode
	Export_pushErr     = pushError
	Export_newRobustDB = func(store db.KvDB) *robustDB {
		return &robustDB{store: store}
	}
)

type expKeyTag struct {
	KeyTag
}

func (ekt *expKeyTag) Equal(want KeyTag) bool {
	return ekt.equal(want)
}
