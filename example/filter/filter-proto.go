package filter

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
)

func UnmarshalProtoFolder(folder string, m map[string]int) {
	filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".js" {
			unmarshalProtoJsfile(path, m)
		}
		return nil
	})

}

func unmarshalProtoJsfile(path string, m map[string]int) {
	///
	key := []byte("addCustomListener(")
	keylen := len(key)
	bs, _ := os.ReadFile(path)
	for {
		idx := bytes.Index(bs, key)
		if idx == -1 {
			break
		}
		bs = bs[idx+keylen:]
		idx = bytes.IndexByte(bs, ',')
		proto := bs[:idx]
		proto = bytes.TrimPrefix(proto, []byte("P.Type"))
		proto = bytes.Trim(proto, "\"")
		m[string(proto)] = 1
	}
}
