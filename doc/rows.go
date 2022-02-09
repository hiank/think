package doc

import (
	"fmt"
	"reflect"
	"strconv"

	"k8s.io/klog/v2"
)

type rowsDoc struct {
	head   []string
	rows   [][]string
	reader RowsReader
}

func (rd *rowsDoc) Encode(v interface{}) (err error) {
	rows, ok := v.([][]string)
	if !ok {
		buf, ok := v.([]byte)
		if !ok || rd.reader == nil {
			return fmt.Errorf("invalid param (only support []byte/[][]string) or no RowsReader (for read []byte to [][]string)")
		}
		if rows, err = rd.reader.Read(buf); err != nil {
			return
		}
	}
	if len(rows) < 2 {
		return fmt.Errorf("at least 2 len needed for rows")
	}
	rd.head, rd.rows = rows[0], rows[1:]
	return
}

func (rd *rowsDoc) Decode(out interface{}) (err error) {
	if m, ok := out.(map[string]interface{}); ok {
		err = rd.decodeM(m)
	} else if l, ok := out.(*[]interface{}); ok {
		err = rd.decodeL(l)
	} else {
		err = fmt.Errorf("invalid param: only support *[]interface{} or map[string]interface{}")
	}
	return
}

func (rd *rowsDoc) Val() []byte {
	return []byte("not support")
}

//typeOf struct type
func (rd *rowsDoc) typeOf(v interface{}) reflect.Type {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv.Type()
}

func (rd *rowsDoc) rangeVal(ftoh map[int]int, rt reflect.Type, f func([]string, interface{})) {
	for _, row := range rd.rows {
		v := reflect.New(rt)
		fv := v.Elem()
		for fidx, hidx := range ftoh {
			if err := decodeToValue(row[hidx], fv.Field(fidx)); err != nil {
				klog.Warning("cannot decode %s to field %s", fv.Field(fidx).Type().Name())
			}
		}
		f(row, v.Interface())
	}
}

func (rd *rowsDoc) decodeL(l *[]interface{}) error {
	if len(*l) != 1 {
		return fmt.Errorf("invalid param: non template in []interface{}")
	}
	rt, out := rd.typeOf((*l)[0]), *l
	if cap(out) < len(rd.rows) {
		out = make([]interface{}, 0, len(rd.rows))
	}
	ftoh, _ := rd.headIndexInRow(rt)
	rd.rangeVal(ftoh, rt, func(row []string, v interface{}) {
		out = append(out, v)
	})
	*l = out
	return nil
}

func (rd *rowsDoc) decodeM(m map[string]interface{}) error {
	if len(m) != 1 {
		return fmt.Errorf("invalid param: non template in map[string]interface{}")
	}
	var ktag string
	var rt reflect.Type
	for k, v := range m {
		ktag, rt = k, rd.typeOf(v)
		delete(m, k) //NOTE: delete template k-v
		break
	}
	ftoh, tinf := rd.headIndexInRow(rt, ktag)
	if tinf == -1 {
		return fmt.Errorf("cannot find Tag(%s) in given tmplate", ktag)
	}
	rd.rangeVal(ftoh, rt, func(row []string, v interface{}) {
		m[row[tinf]] = v
	})
	return nil
}

func (rd *rowsDoc) headIndexInRow(t reflect.Type, ktag ...string) (ftoh map[int]int, tinh int) {
	ftoh, tinh = make(map[int]int), -1
L:
	for fidx := 0; fidx < t.NumField(); fidx++ {
		field := t.Field(fidx)
		tag := field.Tag.Get("excel")
		if tag == "" {
			tag = field.Name
		}
		for hidx, key := range rd.head {
			if key == tag {
				ftoh[fidx] = hidx
				if len(ktag) > 0 && ktag[0] == field.Name {
					tinh = hidx
				}
				continue L
			}
		}
	}
	return
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
