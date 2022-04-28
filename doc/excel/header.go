package excel

import (
	"reflect"
	"strconv"

	// "github.com/hiank/think/doc/v2"
	"github.com/hiank/think/run"
	"golang.org/x/exp/slices"
)

var (
	ErrValueOverflow = run.Err("excel: value overflow")
	ErrNotStructype  = run.Err("excel: not struct type")
)

type Header[T any] struct {
	htof map[int]string
}

func NewHeader[T any](head []string) (*Header[T], error) {
	var tv T
	rt := reflect.TypeOf(tv)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return nil, ErrNotStructype
	}
	htof := make(map[int]string)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get("excel")
		if tag == "" {
			tag = field.Name
		}
		hidx := slices.Index(head, tag)
		if hidx != -1 {
			htof[hidx] = field.Name
		}
	}
	return &Header[T]{htof: htof}, nil
}

//KeyIndex key index in row. for map[string]T 's key value
//return -1 when cannot found it
func (h *Header[T]) KeyIndex(fieldName string) (hidx int) {
	hidx = -1
	for i, s := range h.htof {
		if fieldName == s {
			hidx = i
			break
		}
	}
	return
}

//NewT unmarshal row to T type value
func (h *Header[T]) NewT(row []string) (out T, err error) {
	var rv reflect.Value
	if rt := reflect.TypeOf(out); rt.Kind() == reflect.Ptr {
		out = reflect.New(rt.Elem()).Interface().(T)
		rv = reflect.ValueOf(out)
	} else {
		rv = reflect.ValueOf(&out)
	}
	if rv = rv.Elem(); rv.Kind() != reflect.Struct {
		return out, ErrNotStructype
	}
	for hidx, fname := range h.htof {
		if hidx < len(row) {
			if err = h.unmarshalTo(row[hidx], rv.FieldByName(fname)); err != nil {
				// klog.Warning("doc: cannot decode %s to field %s", rv.FieldByName(fname).Type().Name())
				break
			}
		}
	}
	return
}

func (*Header[T]) unmarshalTo(strv string, out reflect.Value) (err error) {
	switch out.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		if n, err = strconv.ParseInt(strv, 10, 64); err == nil {
			err = limitSet(n, out.OverflowInt, out.SetInt)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var n uint64
		if n, err = strconv.ParseUint(strv, 10, 64); err == nil {
			err = limitSet(n, out.OverflowUint, out.SetUint)
		}
	case reflect.Float32, reflect.Float64:
		var n float64
		if n, err = strconv.ParseFloat(strv, out.Type().Bits()); err == nil {
			err = limitSet(n, out.OverflowFloat, out.SetFloat)
		}
	case reflect.Bool:
		var n bool
		if n, err = strconv.ParseBool(strv); err == nil {
			out.SetBool(n)
		}
	case reflect.String:
		out.SetString(strv)
	}
	return
}

//limitSet set value by `set` with `overflow` limit
func limitSet[T any](v T, overflow func(T) bool, set func(T)) (err error) {
	if !overflow(v) {
		set(v)
	} else {
		err = ErrValueOverflow
	}
	return
}
