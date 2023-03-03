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
	ErrNotProtoMessage = run.Err("doc: not protobuf message")
	ErrNotBytesCoder   = run.Err("doc: coder's base type is not []byte")
)

// Json Coder
type Json []byte

func (j *Json) Decode(v any) (err error) {
	return json.Unmarshal(j.Bytes(), v)
}

func (j *Json) Encode(v any) (err error) {
	switch msg := v.(type) {
	case []byte:
		*j = msg
	default:
		*j, err = json.Marshal(v)
	}
	return
}

func (j *Json) Bytes() []byte {
	return *j
}

func (j *Json) Format() Format {
	return FormatJson
}

// Yaml Coder
type Yaml []byte

func (y *Yaml) Decode(v any) (err error) {
	return yaml.Unmarshal(y.Bytes(), v)
}

func (y *Yaml) Encode(v any) (err error) {
	switch msg := v.(type) {
	case []byte:
		*y = msg
	default:
		*y, err = yaml.Marshal(v)
	}
	return
}

func (y *Yaml) Bytes() []byte {
	return *y
}

func (y *Yaml) Format() Format {
	return FormatYaml
}

// Gob Coder
type Gob []byte

func (g *Gob) Decode(v any) (err error) {
	return gob.NewDecoder(bytes.NewBuffer(g.Bytes())).Decode(v)
}

func (g *Gob) Encode(v any) (err error) {
	if bts, ok := v.([]byte); ok {
		*g = bts
		return
	}
	var b bytes.Buffer
	if err = gob.NewEncoder(&b).Encode(v); err == nil {
		*g = b.Bytes()
	}
	return
}

func (g *Gob) Bytes() []byte {
	return *g
}

func (g *Gob) Format() Format {
	return FormatGob
}

type Proto []byte

func (p *Proto) Decode(v any) (err error) {
	if msg, ok := v.(proto.Message); ok {
		err = proto.Unmarshal(p.Bytes(), msg)
	} else {
		err = ErrNotProtoMessage
	}
	return
}

func (p *Proto) Encode(v any) (err error) {
	switch msg := v.(type) {
	case []byte:
		*p = msg
	case proto.Message:
		*p, err = proto.Marshal(msg)
	default:
		err = ErrNotProtoMessage
	}
	return
}

func (p *Proto) Bytes() []byte {
	return *p
}

func (p *Proto) Format() Format {
	return FormatProto
}

func newBytes[T interface{ Json | Yaml | Gob | Proto }](data []byte) *T {
	var t T = data
	return &t
}

func NewBytesCoder(data []byte, f Format) (c Coder, err error) {
	switch f {
	case FormatGob:
		c = newBytes[Gob](data)
	case FormatJson:
		c = newBytes[Json](data)
	case FormatProto:
		c = newBytes[Proto](data)
	case FormatYaml:
		c = newBytes[Yaml](data)
	default:
		err = ErrUnsupportFormat
	}
	return
}
