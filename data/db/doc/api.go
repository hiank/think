package doc

type Doc interface {
	Decode(v interface{}) error
	Encode(v interface{}) error
	Val() string
}

type Maker interface {
	Make([]byte) Doc
}

var (
	PBMaker   Maker = funcMaker(newPB)
	JsonMaker Maker = funcMaker(newJson)
	GobMaker  Maker = funcMaker(newGob)
)

type funcMaker func([]byte) Doc

func (fm funcMaker) Make(v []byte) Doc {
	return fm(v)
}
