package doc

type Doc interface {
	Decode(out interface{}) error
	Encode(v interface{}) error
	Val() []byte
}

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
