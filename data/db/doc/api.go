package doc

type Doc interface {
	Decode(v interface{}) error
	Encode(v interface{}) error
	Val() string
}
