package excel

import (
	"path/filepath"
	"strings"

	"github.com/hiank/think/run"
	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/slices"
)

const (
	ErrUnsupportFileExt = run.Err("excel: only support file with xlsx/xlsm/xltm/xltx ext")
	ErrFailedToReadRows = run.Err("excel: failed to read rows")
)

var supportedExts = []string{".xlsx", ".xlsm", ".xltm", ".xltx"}

func ReadFileNewMap[KT comparable, VT any](path, ktag string) (m map[KT]VT, err error) {
	rows, err := readFileToRows(path)
	if err == nil {
		m, err = UnmarshalNewMap[KT, VT](rows, ktag)
	}
	return
}

func ReadFileNewSlice[VT any](path string) (s []VT, err error) {
	rows, err := readFileToRows(path)
	if err == nil {
		s, err = UnmarshalNewSlice[VT](rows)
	}
	return
}

func readFileToRows(path string) (rows [][]string, err error) {
	ext := strings.ToLower(filepath.Ext(path))
	if !slices.Contains(supportedExts, ext) {
		return nil, ErrUnsupportFileExt
	}
	////
	var f *excelize.File
	if f, err = excelize.OpenFile(path); err == nil {
		defer f.Close()
		for _, sheetName := range f.GetSheetList() {
			if rows, err = f.GetRows(sheetName); err == nil && len(rows) > 1 {
				///first row storage titles, so the value row start from index 1
				return
			}
		}
		rows, err = nil, ErrFailedToReadRows
	}
	return
}
