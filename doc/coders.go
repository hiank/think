package doc

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/hiank/think/doc/excel"
	"github.com/hiank/think/run"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

var (
	ErrNotProtoMessage  = run.Err("doc: not proto.Message")
	ErrUnsupportType    = run.Err("doc: unsupport param type")
	ErrNotImplemented   = run.Err("doc: not implemented")
	ErrInvalidParamType = run.Err("doc: invalid param type")
)

//internalLimit limited use NewCoder
type internalLimit struct{}

func (internalLimit) internalOnly() {}

//JsonCoder decode/encode json value
type JsonCoder struct{ internalLimit }

func (JsonCoder) Decode(data []byte, out any) error {
	return json.Unmarshal(data, out)
}

func (JsonCoder) Encode(val any) ([]byte, error) {
	return json.Marshal(val)
}

//YamlCoder decode/encode yaml value
type YamlCoder struct{ internalLimit }

func (YamlCoder) Decode(data []byte, out any) error {
	return yaml.Unmarshal(data, out)
}

func (YamlCoder) Encode(val any) ([]byte, error) {
	return yaml.Marshal(val)
}

//GobCoder decode/encode gob value
type GobCoder struct{ internalLimit }

func (GobCoder) Decode(data []byte, out any) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(out)
}

func (GobCoder) Encode(val any) (out []byte, err error) {
	buffer := new(bytes.Buffer)
	if err = gob.NewEncoder(buffer).Encode(val); err == nil {
		out = buffer.Bytes()
	}
	return
}

//PBCoder decode/encode protobuf value
type PBCoder struct{ internalLimit }

func (PBCoder) Decode(data []byte, out any) (err error) {
	msg, ok := out.(proto.Message)
	if ok {
		err = proto.Unmarshal(data, msg)
	} else {
		err = ErrNotProtoMessage
	}
	return
}

func (PBCoder) Encode(val any) (out []byte, err error) {
	msg, ok := val.(proto.Message)
	if ok {
		out, err = proto.Marshal(msg)
	} else {
		err = ErrNotProtoMessage
	}
	return
}

//RowsCoder decode/encode data to map[string]T/*[]T
type RowsCoder struct {
	decoder   excel.Decoder
	unmarshal func(excel.Rows, any) error
}

func (coder *RowsCoder) Decode(data []byte, out any) (err error) {
	rows, err := coder.decoder.UnmarshalNew(data)
	if err == nil {
		err = coder.unmarshal(rows, out)
	}
	return
}

func (*RowsCoder) Encode(any) ([]byte, error) {
	return nil, ErrNotImplemented
}

//NewCoder[T] create given type Coder
func NewCoder[T internalCoder]() Coder {
	var v T
	return v
}

//NewRowsCoder new RowsCoder
//VT is the value type for excel
//exp: want decode to map[string]tmpExcel, VT is tmpExcel
func NewRowsCoder[VT any](kt string, decoders ...excel.Decoder) Coder {
	rc := &RowsCoder{
		unmarshal: func(r excel.Rows, a any) (err error) {
			switch v := a.(type) {
			case map[string]VT:
				err = excel.UnmarshaltoMap(r, v, kt)
			case *[]VT:
				err = excel.UnmarshaltoSlice(r, v)
			default:
				err = ErrInvalidParamType
			}
			return
		},
		decoder: excel.DefaultDecoder,
	}
	if len(decoders) > 0 {
		rc.decoder = decoders[0]
	}
	return rc
}
