package doc

import "reflect"

//Doc[T]
type Doc[T any] struct {
	coder Coder
	tval  T
	bval  []byte
}

//T get saved T type value
func (d *Doc[T]) T() T {
	return d.tval
}

//Bytes get saved []byte type value
func (d *Doc[T]) Bytes() []byte {
	return d.bval
}

//DecodeNew decode []byte to an new T type value. and save them for get use T() Bytes() method
func (d *Doc[T]) DecodeNew(data []byte) (out T, err error) {
	out, err = MakeT[T]()
	if err == nil {
		var val any = out
		if reflect.TypeOf(val).Kind() != reflect.Ptr {
			val = &out
		}
		if err = d.coder.Decode(data, val); err == nil {
			d.bval, d.tval = data, out
		}
	}
	return
}

//Decode decode []byte to passed T type value. and save them for get use T() Bytes() method
func (d *Doc[T]) Decode(data []byte, out T) (err error) {
	if err = d.coder.Decode(data, out); err == nil {
		d.bval, d.tval = data, out
	}
	return
}

//Encode encode T type value to []byte. and save them for get use T() Bytes() method
func (d *Doc[T]) Encode(val T) (out []byte, err error) {
	if out, err = d.coder.Encode(val); err == nil {
		d.bval, d.tval = out, val
	}
	return
}

//New[T] new Doc[T]
func New[T any](coder Coder) *Doc[T] {
	return &Doc[T]{
		coder: coder,
	}
}
