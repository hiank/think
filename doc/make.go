package doc

import "reflect"

//MakeT make a T value
//T cannot be a ** | *map | *slice | chan value
//@param size for make slice. len(size) must be 2
func MakeT[T any](size ...int) (out T, err error) {
	rt := reflect.TypeOf(out)
	switch rt.Kind() {
	case reflect.Ptr:
		rt = rt.Elem()
		switch rt.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice:
			err = ErrUnsupportType
		default:
			out = reflect.New(rt).Interface().(T)
		}
	case reflect.Map:
		out = makeMap(rt, size...).(T)
	case reflect.Slice:
		out = makeSlice(rt, size...).(T)
	case reflect.Chan:
		err = ErrUnsupportType
	}
	return
}

func makeMap(rt reflect.Type, size ...int) any {
	var rv reflect.Value
	if len(size) > 0 {
		rv = reflect.MakeMapWithSize(rt, size[0])
	} else {
		rv = reflect.MakeMap(rt)
	}
	return rv.Interface()
}

func makeSlice(rt reflect.Type, size ...int) any {
	var lenv, capv int
	if len(size) > 1 {
		lenv, capv = size[0], size[1]
	} else if len(size) > 0 {
		lenv, capv = size[0], size[0]
	}
	return reflect.MakeSlice(rt, lenv, capv).Interface()
}
