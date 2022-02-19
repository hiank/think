package doc

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/hiank/think/run"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const (
	ErrNotProtoMessage = run.Err("doc: data must be a proto.Message")
	ErrValueMustBeT    = run.Err("doc: value for encode by docCoder must be a T")
	ErrNilValue        = run.Err("doc: value is nil")
)

// jsonCoder for encode/decode between json format []byte and struct
// decode []byte -> struct
// encode struct -> []byte
type jsonCoder struct{}

func (jsonCoder) Decode(data []byte, out interface{}) error {
	return json.Unmarshal(data, out)
}

func (jsonCoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// yamlCoder for encode/decode between yaml format []byte and struct
// decode []byte -> struct
// encode struct -> []byte
type yamlCoder struct{}

func (yamlCoder) Decode(data []byte, out interface{}) error {
	return yaml.Unmarshal(data, out)
}

func (yamlCoder) Encode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// gobCoder for encode/decode between gob format []byte and struct
// decode []byte -> struct
// encode struct -> []byte
type gobCoder struct{}

func (gobCoder) Decode(data []byte, out interface{}) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(out)
}

func (gobCoder) Encode(v interface{}) (out []byte, err error) {
	buf := new(bytes.Buffer)
	if err = gob.NewEncoder(buf).Encode(v); err == nil {
		out = buf.Bytes()
	}
	return
}

// protoCoder for encode/decode between proto.Message format []byte and struct
// decode []byte -> struct
// encode struct -> []byte
type protoCoder struct{}

func (protoCoder) Decode(data []byte, out interface{}) error {
	if msg, ok := out.(proto.Message); ok {
		return proto.Unmarshal(data, msg)
	}
	return ErrNotProtoMessage
}

func (protoCoder) Encode(v interface{}) ([]byte, error) {
	if msg, ok := v.(proto.Message); ok {
		return proto.Marshal(msg)
	}
	return nil, ErrNotProtoMessage
}

// Tcoder convert T to Coder
// encode T's V to out []byte
// decode []byte to T's V
type Tcoder byte

//Encode encode v to bytes
//NOTE: v must be doc.Doc
func (dc Tcoder) Encode(v interface{}) (out []byte, err error) {
	t, err := dc.check(v)
	if err == nil {
		out, err = t.Encode()
	}
	return
}

func (dc Tcoder) Decode(data []byte, out interface{}) (err error) {
	t, err := dc.check(out)
	if err == nil {
		err = t.Decode(data)
	}
	return
}

func (Tcoder) check(v interface{}) (t T, err error) {
	t, ok := v.(T)
	if !ok {
		err = ErrValueMustBeT
	} else if t.V == nil {
		err = ErrNilValue
	}
	return
}
