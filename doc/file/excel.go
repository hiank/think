package file

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
	"k8s.io/klog/v2"
)

//Excel reader for excel file
//	 - ".xlsx"
//   - ".xlsm"
//   - ".xltm"
//   - ".xltx"
//NOTE: only support first not empty sheet
type excelRowsReader byte

func (reader excelRowsReader) Read(v []byte) (rows [][]string, err error) {
	f, err := excelize.OpenReader(bytes.NewReader(v))
	if err != nil {
		klog.Warning(err)
		return
	}
	defer f.Close()

	for _, sheetName := range f.GetSheetList() {
		if rows, err = f.GetRows(sheetName); err == nil && len(rows) > 1 {
			///first row storage titles, so the value row start from index 1
			return
		}
	}
	return nil, fmt.Errorf("invalid param: cannot read legitimate rows value")
}
