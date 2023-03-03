package excel

import "reflect"

// var (
// 	Export_rangeRows[T any] = rangeRows[T]
// )

// func Export_rangeRows[T any](rows [][]string, call func(*Header[T], []string) error) error {
// 	return nil //rangeRows(rows, call)
// }

var (
	Export_readFileToRows = readFileToRows
	Export_newTitle       = newTitle
	Export_convertoValue  = convertoValue
	// Export_newRows[T]        = newRows[T]
)

func Export_newRows[T any](raw [][]string) (*Rows, error) {
	return newRows[T](raw)
}

func ExportGetTitleMember(t *title) ([]string, map[int]string) {
	return t.row, t.m
}

func ExportCallRowsNew(r *Rows, row []string) reflect.Value {
	return r.new(row)
}
