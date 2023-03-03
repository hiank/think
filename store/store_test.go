package store_test

import (
	"fmt"
	"testing"

	"github.com/hiank/think/store"
	"gotest.tools/v3/assert"
)

type tmpKey string

type tmpDictionary struct {
}

func (td *tmpDictionary) Scan(k tmpKey, out any) (suc bool, err error) {
	if v, ok := out.(*tmpKey); ok {
		*v = k
		suc = true
	} else {
		err = errInvalidScan
	}
	return
}

func (td *tmpDictionary) Set(k tmpKey, v any) error {
	if tv, ok := v.(tmpKey); ok {
		if tv == k {
			return nil
		}
	}
	return errInvalidSet
}

func (td *tmpDictionary) Del(k tmpKey, outs ...any) error {
	return errUnimplemented
}

func (td *tmpDictionary) Close() error {
	return errUnimplemented
}

var (
	errUnimplemented = fmt.Errorf("unimplemented")
	errInvalidSet    = fmt.Errorf("invalid set")
	errInvalidScan   = fmt.Errorf("invalid scan")
)

func TestDictionary(t *testing.T) {
	var d store.Dictionary[tmpKey] = &tmpDictionary{}
	d.Scan(tmpKey(""), nil)
	// d =
}

func TestJsonkey(t *testing.T) {
	var jk store.Jsonkey

	v, found := jk.Get("coll")
	assert.Equal(t, found, false)
	assert.Equal(t, v, "")
	(&jk).Encode(store.JsonkeyPair{K: "coll", V: "1"})
	v, found = jk.Get("coll")
	// assert.Equal(t, err, nil, err)
	assert.Assert(t, found)
	assert.Equal(t, v, "1")

	v, found = jk.Get("doc")
	assert.Equal(t, found, false)
	assert.Equal(t, v, "")
}

func TestConvertoEasy(t *testing.T) {
	var d tmpDictionary
	ed := store.ConvertoEasy[tmpKey](&d)
	var out tmpKey
	found, err := ed.Scan("fm", &out)
	assert.Assert(t, found)
	assert.Equal(t, err, nil, err)
	assert.Equal(t, out, tmpKey("fm"))

	var out2 string
	found, err = ed.Scan("fm2", &out2)
	assert.Assert(t, !found)
	assert.Equal(t, err, errInvalidScan)

	assert.Equal(t, ed.Set("fm", tmpKey("fm")), nil, err)
	assert.Equal(t, ed.Set("fm", tmpKey("fm2")), errInvalidSet)

	assert.Equal(t, ed.Del("fm"), errUnimplemented)
	assert.Equal(t, ed.Close(), errUnimplemented)
}
