package sys

import (
	"github.com/hiank/think/doc"
	"github.com/hiank/think/run"
)

var (
	ErrUnsupportFormat = run.Err("sys: unsupport format (currently support 'json'/'yaml')")
)

type Bytes struct {
	buffer []byte
	doc.Coder
}

func (b *Bytes) UnmarshalTo(out any) error {
	return b.Decode(b.buffer, out)
}

func formatoBytes(f Format, get func() ([]byte, error)) (b *Bytes, err error) {
	var coder doc.Coder
	switch f {
	case FormatJson:
		coder = doc.NewCoder[doc.JsonCoder]()
	case FormatYaml:
		coder = doc.NewCoder[doc.YamlCoder]()
	default:
		return nil, ErrUnsupportFormat
	}
	buffer, err := get()
	if err == nil {
		b = &Bytes{buffer: buffer, Coder: coder}
	}
	return
}
