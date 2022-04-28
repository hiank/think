package excel

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/hiank/think/run"
	"github.com/xuri/excelize/v2"
)

var (
	ErrFailedReadRows  = run.Err("excel: failed to read rows")
	ErrUnsupportSuffix = run.Err("excel: unsupport file suffix (must one of xlsx/xlsm/xltm/xltx)")
)

//DefaultDecoder us excelize pacage
var DefaultDecoder Decoder = defaultDecoder{}

//Excel reader for excel file
//	 - ".xlsx"
//   - ".xlsm"
//   - ".xltm"
//   - ".xltx"
//NOTE: only support first not empty sheet
type defaultDecoder struct{}

//UnmarshalNew unmarshal []byte to Rows
func (defaultDecoder) UnmarshalNew(v []byte) (rows Rows, err error) {
	f, err := excelize.OpenReader(bytes.NewReader(v))
	if err != nil {
		// klog.Warning(err)
		return
	}
	defer f.Close()

	for _, sheetName := range f.GetSheetList() {
		if rows, err = f.GetRows(sheetName); err == nil && len(rows) > 1 {
			///first row storage titles, so the value row start from index 1
			return
		}
	}
	return nil, ErrFailedReadRows
}

func FiletoMap[VT any](path string, kt string) (m map[string]VT, err error) {
	rows, err := FiletoRows(path)
	if err == nil {
		m = make(map[string]VT)
		err = UnmarshaltoMap(rows, m, kt)
	}
	return
}

func FiletoSlice[VT any](path string) (s []VT, err error) {
	rows, err := FiletoRows(path)
	if err == nil {
		s = make([]VT, 0, len(rows))
		err = UnmarshaltoSlice(rows, &s)
	}
	return
}

func FiletoRows(path string) (rows Rows, err error) {
	data, err := readExcelFile(path)
	if err == nil {
		rows, err = DefaultDecoder.UnmarshalNew(data)
	}
	return
}

func readExcelFile(path string) (data []byte, err error) {
	err = ErrUnsupportSuffix
	if idx := strings.LastIndexByte(path, '.'); idx != -1 {
		switch strings.ToLower(path[idx+1:]) {
		case "xlsx", "xlsm", "xltm", "xltx":
			data, err = ioutil.ReadFile(path)
		}
	}
	return
}
