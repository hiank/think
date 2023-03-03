package excel

import (
	"reflect"
	"strconv"

	"github.com/hiank/think/run"
	"golang.org/x/exp/slices"
	"k8s.io/klog/v2"
)

var (
	ErrInvalidParam   = run.Err("excel: invalid param")
	ErrNonKeyField    = run.Err("excel: non key field in V for map")
	ErrInvalidKeyType = run.Err("excel: invalid key type for map")
	ErrOverflow       = run.Err("excel: value overflow")
	ErrNotStruct      = run.Err("excel: not type struct")
)

func UnmarshalNewSlice[T any](rows [][]string) (s []T, err error) {
	r, err := newRows[T](rows)
	if err != nil {
		return
	}
	s = make([]T, 0, len(rows)-1)
	err = r.RangeNew(func(v reflect.Value) error {
		s = append(s, v.Interface().(T))
		return nil
	})
	return
}

func UnmarshalNewMap[KT comparable, VT any](rows [][]string, ktag string) (m map[KT]VT, err error) {
	r, err := newRows[VT](rows)
	if err != nil {
		return
	}
	////
	fn, found := r.GetFieldName(ktag)
	if !found {
		return m, ErrNonKeyField
	}
	////
	var kv KT
	krt, m := reflect.TypeOf(kv), map[KT]VT{}
	err = r.RangeNew(func(v reflect.Value) error {
		elem := v
		if elem.Kind() == reflect.Pointer {
			elem = elem.Elem()
		}
		rv := elem.FieldByName(fn)
		if !rv.CanConvert(krt) {
			return ErrInvalidKeyType
		}
		key := rv.Convert(krt).Interface().(KT)
		m[key] = v.Interface().(VT)
		return nil
	})
	return
}

type title struct {
	////
	row []string
	m   map[int]string //map[title row index]field name
}

func newTitle(row []string, vrt reflect.Type) (t *title, err error) {
	rt := vrt
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	////
	t = &title{row: row, m: map[int]string{}}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get("excel")
		if tag == "" {
			tag = field.Name
		}
		hidx := slices.Index(row, tag)
		if hidx != -1 {
			t.m[hidx] = field.Name
		}
	}
	return
}

// limitSet set value by `set` with `overflow` limit
func limitSet[T any](v T, overflow func(T) bool, set func(T)) (err error) {
	if !overflow(v) {
		set(v)
	} else {
		err = ErrOverflow
	}
	return
}

// /
func convertoValue(strv string, out reflect.Value) (err error) {
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

type Rows struct {
	///
	raw [][]string //raw data
	rt  reflect.Type
	*title
}

// newRows
// @param rt : value type
func newRows[T any](raw [][]string) (r *Rows, err error) {
	if len(raw) < 2 {
		return r, ErrInvalidParam
	}
	var v T
	rt := reflect.TypeOf(v)
	t, err := newTitle(raw[0], rt)
	if err == nil {
		r = &Rows{raw: raw, rt: rt, title: t}
	}
	return
}

func (r *Rows) GetFieldName(tag string) (fn string, found bool) {
	rt := r.rt
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		ftag := field.Tag.Get("excel")
		if ftag == "" {
			ftag = field.Name
		}
		if found = (tag == ftag); found {
			fn = field.Name
			break
		}
	}
	return
}

func (r *Rows) RangeNew(exec func(reflect.Value) error) (err error) {
	///
	for _, row := range r.raw[1:] {
		if err = exec(r.new(row)); err != nil {
			break
		}
	}
	return
}

func (r *Rows) new(row []string) (v reflect.Value) {
	var elemv reflect.Value
	if r.rt.Kind() == reflect.Pointer {
		v = reflect.New(r.rt.Elem())
		elemv = v.Elem()
	} else {
		v = reflect.New(r.rt).Elem()
		elemv = v
	}
	/////
	max := len(row)
	for i, fn := range r.m {
		if i < max {
			if e := convertoValue(row[i], elemv.FieldByName(fn)); e != nil {
				klog.Warning("excel: convert string to Field data error:", e)
			}
		}
	}
	return
}
