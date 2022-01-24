package fp

import (
	"errors"
	"reflect"
	"strconv"
)

type rowsConv struct {
	rows [][]string
	t    reflect.Type
	finr []int //[fieldIndex]rowIndex
}

func newRowsConv(rows [][]string, t reflect.Type) *rowsConv {
	rc := &rowsConv{rows: rows[1:], t: t}
	rc.initFindexInRow(rows[0])
	return rc
}

func (rc *rowsConv) Unmarshal() (out []interface{}) {
	out = make([]interface{}, len(rc.rows))
	for i, row := range rc.rows {
		v := reflect.New(rc.t) //makeVal()
		fv := v.Elem()
		for fi, ri := range rc.finr {
			if ri != -1 {
				rc.decodeToValue(row[ri], fv.Field(fi))
			}
		}
		out[i] = v.Interface()
	}
	return
}

func (rc *rowsConv) decodeToValue(strv string, v reflect.Value) (err error) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		if n, err = strconv.ParseInt(strv, 10, 64); err == nil {
			if !v.OverflowInt(n) {
				v.SetInt(n)
			} else {
				err = errors.New("overflow for int64")
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var n uint64
		if n, err = strconv.ParseUint(strv, 10, 64); err == nil {
			if !v.OverflowUint(n) {
				v.SetUint(n)
			} else {
				err = errors.New("overflow for uint64")
			}
		}
	case reflect.Float32, reflect.Float64:
		var n float64
		if n, err = strconv.ParseFloat(strv, v.Type().Bits()); err == nil {
			if !v.OverflowFloat(n) {
				v.SetFloat(n)
			} else {
				err = errors.New("overflow for float64")
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

func (rc *rowsConv) initFindexInRow(keys []string) {
	finr := make([]int, rc.t.NumField())
L:
	for fi := range finr {
		sf := rc.t.Field(fi)
		tag := sf.Tag.Get("excel")
		if tag == "" {
			tag = sf.Name
		}
		for ki, key := range keys {
			if key == tag {
				finr[fi] = ki
				continue L
			}
		}
		finr[fi] = -1
	}
	rc.finr = finr
}
