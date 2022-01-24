package doc

import "encoding/json"

type Json []byte

func newJson(v []byte) Doc {
	js := new(Json)
	if v != nil {
		*js = v
	}
	return js
}

func (js *Json) Decode(v interface{}) error {
	return json.Unmarshal(*js, v)
}

func (js *Json) Encode(v interface{}) error {
	buf, err := json.Marshal(v)
	if err == nil {
		*js = Json(buf)
	}
	return err
}

func (js *Json) Val() string {
	return string(*js)
}
