package db

import (
	"fmt"
	"io"
	"regexp"
)

//KvDB key-value stype database client
type KvDB interface {
	Get(k string, out interface{}) (found bool, err error)
	Set(k string, v interface{}) error
	Delete(k string) error
	io.Closer
}

const (
	//KTMem tag for use memory store
	KTMem KeyTag = 1 << 0
	//KTDisk tag for use disk store
	KTDisk KeyTag = 1 << 1
	//KTMix tag for mix mode
	//use both memory and disk store
	//NOTE: using mixed mode is prone to data inconsistency and should be avoided
	KTMix KeyTag = KTMem | KTDisk

	ktRule   string = "[%d@KT]%s"
	ktRegexp string = `^\[(.*)@KT\]`
)

type KeyTag uint8

//Encode encode given baseKey to the key contains tag value
func (kt KeyTag) Encode(baseKey string) (key string) {
	r := regexp.MustCompile(ktRegexp)
	if loc := r.FindStringIndex(baseKey); len(loc) > 0 {
		baseKey = baseKey[loc[1]:]
	}
	return fmt.Sprintf(ktRule, kt, baseKey)
}

//equal check if the kt contians want tag
func (kt KeyTag) equal(want KeyTag) bool {
	return (want > 0) && (kt&want) == want
}

//DBS database set
//temporarily only supports key-value database
type DBS interface {
	//KvDB key-value database store
	KvDB() KvDB
}
