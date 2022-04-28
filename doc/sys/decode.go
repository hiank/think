package sys

import (
	"io/ioutil"
	"reflect"

	"github.com/hiank/think/doc"
)

//UnmarshalNew unmarshal given filepath file data to new T type object
func UnmarshalNew[T any](filepath string) (out T, err error) {
	if out, err = doc.MakeT[T](); err == nil {
		var v any = out
		if reflect.TypeOf(out).Kind() != reflect.Ptr {
			v = &out
		}
		err = UnmarshalTo(filepath, v)
	}
	return
}

//UnmarshalTo unmarshal given filepath file data to given T type object
func UnmarshalTo[T any](filepath string, out T) (err error) {
	b, err := formatoBytes(formatFromPath(filepath), func() ([]byte, error) {
		return ioutil.ReadFile(filepath)
	})
	if err == nil {
		err = b.UnmarshalTo(out)
	}
	return
}
