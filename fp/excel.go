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
				rc.setFieldValue(fv.Field(fi), row[ri])
			}
		}
		out[i] = v.Interface()
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

func (rc *rowsConv) setFieldValue(v reflect.Value, s string) (err error) {
	overflow := false
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		if n, err = strconv.ParseInt(s, 10, 64); err == nil {
			if overflow = v.OverflowInt(n); !overflow {
				v.SetInt(n)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var n uint64
		if n, err = strconv.ParseUint(s, 10, 64); err == nil {
			if overflow = v.OverflowUint(n); !overflow {
				v.SetUint(n)
			}
		}
	case reflect.Float32, reflect.Float64:
		var n float64
		if n, err = strconv.ParseFloat(s, v.Type().Bits()); err == nil {
			if overflow = v.OverflowFloat(n); !overflow {
				v.SetFloat(n)
			}
		}
	case reflect.Bool:
		var n bool
		if n, err = strconv.ParseBool(s); err == nil {
			v.SetBool(n)
		}
	case reflect.String:
		v.SetString(s)
	}
	if overflow {
		err = errors.New("value in excel is overflow for want type")
	}
	return
}
