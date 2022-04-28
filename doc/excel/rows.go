package excel

import (
	"github.com/hiank/think/run"
)

var (
	ErrInvalidParamType = run.Err("excel: invalid param type")
	ErrNonKeyFound      = run.Err("excel: non key found for map")
)

func rangeRows[T any](rows Rows, call func(*Header[T], []string) error) (err error) {
	if len(rows) < 2 {
		return ErrInvalidParamType
	}
	header, err := NewHeader[T](rows[0])
	if err != nil {
		return
	}
	for _, row := range rows[1:] {
		if err = call(header, row); err != nil {
			break
		}
	}
	return
}

func UnmarshaltoMap[T any](rows Rows, out map[string]T, kt string) error {
	return rangeRows(rows, func(header *Header[T], row []string) (callErr error) {
		idx := header.KeyIndex(kt)
		if idx != -1 && idx < len(row) {
			out[row[idx]], callErr = header.NewT(row)
		} else {
			callErr = ErrNonKeyFound
		}
		return
	})
}

func UnmarshaltoSlice[T any](rows Rows, out *[]T) (err error) {
	slice := make([]T, 0, len(rows))
	if err = rangeRows(rows, func(header *Header[T], row []string) (callErr error) {
		v, callErr := header.NewT(row)
		if callErr == nil {
			slice = append(slice, v)
		}
		return
	}); err == nil {
		*out = slice
	}
	return
}
