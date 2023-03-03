package filter

import (
	"bytes"
	"fmt"
)

func UnmarshalFuncidFolder(folder string, m map[string]int) {
	WalkGivenExts(folder, func(path string) error {
		unmarshalFuncidJsfile(path, m)
		return nil
	}, ".js")
}

func unmarshalFuncidJsfile(path string, m map[string]int) {
	key := []byte("FuncId.")
	keylen := len(key)
	// bs, _ := os.ReadFile(path)
	bs := ReadJs(path)
	for {
		idx := bytes.Index(bs, key)
		if idx == -1 {
			break
		}
		bs = bs[idx+keylen:]
		idx = bytes.IndexFunc(bs, func(r rune) bool {
			return r == ',' || r == ']' || r == '\n'
		})
		// idx = bytes.IndexByte(bs, ',')
		proto := bs[:idx]
		fmt.Println(string(proto))
		m[string(proto)] = 1
	}
}
