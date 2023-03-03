package store

import (
	"encoding/json"
)

type jsonMap struct {
	M map[string]string `json:"jsonkey.map"`
}

type JsonkeyPair struct {
	K, V string
}

type Jsonkey string

//Encode store k-v
//NOTE: unsafe async
func (jk *Jsonkey) Encode(pairs ...JsonkeyPair) {
	jm := &jsonMap{M: make(map[string]string)}
	json.Unmarshal([]byte(*jk), &jm) //ignore unmarshal error here
	for _, pair := range pairs {
		jm.M[pair.K] = pair.V
	}
	b, _ := json.Marshal(&jm) //must not go wrong
	*jk = Jsonkey(b)
}

//Get get v for k
func (jk Jsonkey) Get(k string) (v string, found bool) {
	var jm jsonMap
	if err := json.Unmarshal([]byte(jk), &jm); err == nil {
		v, found = jm.M[k]
	}
	return
}
