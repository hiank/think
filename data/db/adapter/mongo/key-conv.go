package mongo

import "strings"

const (
	defaultCollKey string = "_single_coll"
)

type keyConv struct {
	coll string
	doc  string
}

func newKeyConv(k string) *keyConv {
	k = strings.TrimPrefix(k, "@")
	k = strings.TrimSuffix(k, "@")
	kv := &keyConv{
		coll: defaultCollKey,
		doc:  k,
	}
	if idx := strings.IndexByte(k, '@'); idx != -1 {
		kv.doc, kv.coll = k[:idx], k[idx+1:]
	}
	return kv
}

func (conv *keyConv) GetColl() string {
	return conv.coll
}

func (conv *keyConv) GetDoc() string {
	return conv.doc
}
