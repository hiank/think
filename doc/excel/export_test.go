package excel

// var ( 
// 	Export_rangeRows[T any] = rangeRows[T]
// )

func Export_rangeRows[T any](rows [][]string, call func(*Header[T], []string) error) error {
	return rangeRows(rows, call)
}

var (
	Export_readExcelFile = readExcelFile
)
