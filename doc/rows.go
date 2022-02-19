package doc

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/hiank/think/run"
	"k8s.io/klog/v2"
)

const (
	ErrNotSliceptrOrMap = run.Err("doc: RowsCoder only support slice(ptr when decode)|map")
	ErrNonTemplateValue = run.Err("doc: requires a template value for decode")
	ErrUnimplemented    = run.Err("doc: unimplemented method")
)

type RowsCoder struct {
	// Key Tag for map's key (value of one of rows[0]'s value)
	KT string
	RC RowsConverter
}

func (rd RowsCoder) Decode(data []byte, out interface{}) (err error) {
	rows, err := rd.RC.ToRows(data)
	if err != nil {
		//don't pass here often, so don't need to consider process optimization
		return
	}
	rv := reflect.ValueOf(out)
	switch rv.Kind() {
	case reflect.Map:
		err = rd.toMap(rows, rv)
	case reflect.Ptr: //ptr for slice
		err = rd.toArray(rows, rv.Elem())
	default:
		err = ErrNotSliceptrOrMap
	}
	return
}

func (rd RowsCoder) Encode(v interface{}) (out []byte, err error) {
	// out, ok := v.([]byte)
	// if !ok {
	// 	err = ErrUnimplemented
	// }
	return nil, ErrUnimplemented
}

//parseHead parse head slice
//vt cannot be Ptr
//@return map[fieldindex]sliceindex, key index in field
func (rd RowsCoder) parseHead(h []string, vt reflect.Type) (ftoh map[int]int, kidx int) {
	ftoh, kt := make(map[int]int), rd.KT
	if kt == "" {
		kt = "ID"
	}
L:
	for fidx := 0; fidx < vt.NumField(); fidx++ {
		field := vt.Field(fidx)
		tag := field.Tag.Get("excel")
		if tag == "" {
			tag = field.Name
		}
		for hidx, key := range h {
			if key == tag {
				ftoh[fidx] = hidx
				if kt == field.Name {
					kidx = fidx
				}
				continue L
			}
		}
	}
	return
}

func (rd RowsCoder) rowToValue(row []string, vt reflect.Type, ftoh map[int]int) (v reflect.Value) {
	v = reflect.New(vt)
	fv := v.Elem()
	for fidx, hidx := range ftoh {
		var s string
		if hidx < len(row) { //fill in blanks
			s = row[hidx]
		}
		if err := decodeToValue(s, fv.Field(fidx)); err != nil {
			klog.Warning("doc: cannot decode %s to field %s", fv.Field(fidx).Type().Name())
		}
	}
	return
}

func (rd RowsCoder) toMap(rows [][]string, rv reflect.Value) error {
	rd.rangeVal(rows, rv, func(v reflect.Value, i int) {
		sv := v
		if v.Kind() == reflect.Ptr {
			sv = v.Elem()
		}
		rv.SetMapIndex(sv.Field(i), v)
	})
	return nil
}

func (rd RowsCoder) toArray(rows [][]string, rv reflect.Value) error {
	if rv.Kind() != reflect.Slice {
		return ErrNotSliceptrOrMap
	}
	tmp := rv
	rd.rangeVal(rows, rv, func(v reflect.Value, _ int) {
		tmp = reflect.Append(tmp, v)
	})
	rv.Set(tmp)
	return nil
}

//out must be slice or map
func (rd RowsCoder) rangeVal(rows [][]string, out reflect.Value, f func(reflect.Value, int)) {
	vt := out.Type().Elem() //value of slice or map
	ptr := vt.Kind() == reflect.Ptr
	if ptr {
		if vt = vt.Elem(); vt.Kind() == reflect.Ptr {
			klog.Warning("doc: cannot decode to pointer-pointer-value")
			return
		}
	}
	ftoh, kidx := rd.parseHead(rows[0], vt)
	for _, row := range rows[1:] {
		v := rd.rowToValue(row, vt, ftoh)
		if !ptr {
			v = v.Elem()
		}
		f(v, kidx)
	}
}

func decodeToValue(strv string, v reflect.Value) (err error) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		if n, err = strconv.ParseInt(strv, 10, 64); err == nil {
			if !v.OverflowInt(n) {
				v.SetInt(n)
			} else {
				err = fmt.Errorf("overflow for int64")
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var n uint64
		if n, err = strconv.ParseUint(strv, 10, 64); err == nil {
			if !v.OverflowUint(n) {
				v.SetUint(n)
			} else {
				err = fmt.Errorf("overflow for uint64")
			}
		}
	case reflect.Float32, reflect.Float64:
		var n float64
		if n, err = strconv.ParseFloat(strv, v.Type().Bits()); err == nil {
			if !v.OverflowFloat(n) {
				v.SetFloat(n)
			} else {
				err = fmt.Errorf("overflow for float64")
			}
		}
	case reflect.Bool:
		var n bool
		if n, err = strconv.ParseBool(strv); err == nil {
			v.SetBool(n)
		}
	case reflect.String:
		v.SetString(strv)
	}
	return
}
