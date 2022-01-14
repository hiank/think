package fp

import (
	"reflect"

	"github.com/xuri/excelize/v2"
	"k8s.io/klog/v2"
)

//IParser text parse to code's object
//   - ".json": json, support LoadFile and LoadJsonBytes
//   - ".yaml": yaml, support LoadFile and LoadYamlBytes
//NOTE: not support ".xls" format
type IParser interface {
	LoadFile(paths ...string)
	LoadYamlBytes(values []byte)
	LoadJsonBytes(values []byte)

	//ParseAndClear parse loaded text to given objects
	//the Parser will clear loaded text after parse
	ParseAndClear(configs ...interface{})
}

//IExcel parse value from excel file
//	 - ".xlsx"
//   - ".xlsm"
//   - ".xltm"
//   - ".xltx"
//NOTE: only support first not empty sheet
func UnmarshalExcel(path string, t reflect.Type) (out []interface{}) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		klog.Warning(err)
		return
	}
	defer f.Close()

	for _, sheetName := range f.GetSheetList() {
		if rows, err := f.GetRows(sheetName); err == nil && len(rows) > 1 {
			///first row storage titles, so the value row start from index 1
			out = newRowsConv(rows, t).Unmarshal()
			break
		}
	}
	return
}
