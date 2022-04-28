//Package sys implements method to load config buffer and unmarshal to object
//current support 'json' 'yaml' format data
package sys

import "strings"

type Format int

const (
	FormatUnsupport Format = iota
	FormatJson
	FormatYaml
	// FormatGob
	// FormatPB
	// FormatRows
)

func formatFromPath(path string) (f Format) {
	if idx := strings.LastIndexByte(path, '.'); idx != -1 {
		switch strings.ToLower(path[idx+1:]) {
		case "json":
			f = FormatJson
		case "yaml":
			f = FormatYaml
		}
	}
	return
}
