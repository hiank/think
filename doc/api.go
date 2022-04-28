package doc

type Coder interface {
	Decode(data []byte, out any) error
	Encode(val any) ([]byte, error)
}

//internalCoder limited use NewCoder
type internalCoder interface {
	Coder
	internalOnly()
}
