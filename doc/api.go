package doc

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hiank/think/run"
)

type Format int

const (
	FormatUndefined  Format = 0
	formatStyleBytes Format = 1 << 0
	FormatJson       Format = (1 << 1) | formatStyleBytes //Json format
	FormatYaml       Format = (1 << 2) | formatStyleBytes //Yaml format
	FormatGob        Format = (1 << 3) | formatStyleBytes //Gob format
	FormatProto      Format = (1 << 4) | formatStyleBytes //Proto format
)

var (
	ErrUnsupportFormat = run.Err("doc: unsupport format")
)

type Coder interface {
	Decode(any) error
	Encode(any) error
	Bytes() []byte
	Format() Format
}

// FilepathToFormat filepath ext to Format
func FilepathToFormat(path string) (f Format) {
	ext := filepath.Ext(path)
	switch strings.ToLower(ext) {
	case ".json":
		f = FormatJson
	case ".yaml":
		f = FormatYaml
	case ".gob":
		f = FormatGob
	case ".proto":
		f = FormatProto
	}
	return f
}

// ReadFile to bytes coder (base for filename ext)
func ReadFile(path string) (c Coder, err error) {
	var b []byte
	if b, err = os.ReadFile(path); err == nil {
		c, err = NewBytesCoder(b, FilepathToFormat(path))
	}
	return
}
