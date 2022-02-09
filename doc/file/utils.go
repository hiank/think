package file

import (
	"fmt"
	"strings"
)

func pushError(err, ex error) error {
	if err == nil {
		return ex
	}
	if ex != nil {
		err = fmt.Errorf("%s&&%s", err.Error(), ex.Error())
	}
	return err
}

func pathToForm(path string) (f Form) {
	f = FormInvalid
	if idx := strings.LastIndexByte(path, '.'); idx != -1 {
		switch strings.ToLower(path[idx+1:]) {
		case "json":
			f = FormJson
		case "yaml":
			f = FormYaml
		case "xlsx", "xlsm", "xltm", "xltx":
			f = FormRows
		}
	}
	return
}
