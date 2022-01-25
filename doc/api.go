package doc

type Doc interface {
	Decode(out interface{}) error
	Encode(v interface{}) error
	Val() []byte
}

type BytesMaker interface {
	Make([]byte) Doc
}

var (
	PBMaker   BytesMaker = funcBytesMaker(func(b []byte) Doc { pb := PB(b); return &pb })
	JsonMaker BytesMaker = funcBytesMaker(func(b []byte) Doc { js := Json(b); return &js })
	YamlMaker BytesMaker = funcBytesMaker(func(b []byte) Doc { ym := Yaml(b); return &ym })
	GobMaker  BytesMaker = funcBytesMaker(func(b []byte) Doc { gb := Gob(b); return &gb })
)

//RowsReader read []byte data to [][]string data
type RowsReader interface {
	Read(v []byte) ([][]string, error)
}

//NewRows new rowsDoc
func NewRows(reader RowsReader, dvals ...[][]string) Doc {
	rd := &rowsDoc{reader: reader}
	if len(dvals) > 0 {
		rd.Encode(dvals[0])
	}
	return rd
}
