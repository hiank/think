package excel

type Rows [][]string

type Decoder interface {
	UnmarshalNew([]byte) (Rows, error)
}
