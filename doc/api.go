package doc

//RowsConverter convert between []byte and [][]string
type RowsConverter interface {
	ToRows([]byte) ([][]string, error)
	ToBytes([][]string) ([]byte, error)
}

//Coder decode|encode between bytes and struct
type Coder interface {
	Decode(b []byte, out interface{}) error
	Encode(v interface{}) ([]byte, error)
}

//Maker for make T *B value
type Maker interface {
	MakeT(v interface{}) T
	MakeB(d []byte) *B
}

var (
	// Yaml maker use yamlCoder
	Y Maker = &maker{coder: yamlCoder{}}
	// Json maker use jsonCoder
	J Maker = &maker{coder: jsonCoder{}}
	// Gob maker use gobCoder
	G Maker = &maker{coder: gobCoder{}}
	// Proto maker use protoCoder
	P Maker = &maker{coder: protoCoder{}}
)
