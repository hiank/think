package easy

import (
	"reflect"

	"github.com/hiank/think/run"
)

const (
	// ErrUseMakeInstead = run.Err("easy: use InstantiateMake instead")
	ErrUnsupportType = run.Err("easy: unsupport type for instantiate")
)

// Instantiate struct or struct pointer
// NOTE: no support other type
func Instantiate[T any]() (v T, err error) {
	///
	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.Ptr:
		rt = rt.Elem()
		switch rt.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice:
			err = ErrUnsupportType
		default:
			v = reflect.New(rt).Interface().(T)
		}
	case reflect.Map, reflect.Slice, reflect.Chan:
		err = ErrUnsupportType
	}
	return
}

// func Make[T any](sizes ...int) (v T, err error) {
// 	rt := reflect.TypeOf(v)
// 	switch rt.Kind() {
// 	case reflect.Slice:
// 		params := copyNew(2, sizes)
// 		v = reflect.MakeSlice(rt, params[0], params[1]).Interface().(T)
// 	case reflect.Map:
// 		params := copyNew(1, sizes)
// 		v = reflect.MakeMapWithSize(rt, params[0]).Interface().(T)
// 	case reflect.Chan:
// 		params := copyNew(1, sizes)
// 		v = reflect.MakeChan(rt, params[0]).Interface().(T)
// 	default:
// 		err = ErrUnsupportType
// 	}
// 	return
// }

// func copyNew[T any](len int, src []T) []T {
// 	v := make([]T, len)
// 	copy(v, src)
// 	return v
// }
