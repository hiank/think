package filter

import (
	"bytes"
	"unicode"
)

func ListHUD_LISTUsed(folder string, m map[string]int) {
	WalkGivenExts(folder, func(path string) error {
		unmarshalHudJsfile(path, m)
		return nil
	}, ".js")
}

func unmarshalHudJsfile(path string, m map[string]int) {
	key, suffixkey := []byte("HUD_LIST."), []byte(".getCom(")
	keylen, suffixkeylen := len(key), len(suffixkey)
	bs := ReadJs(path)
	suffixbs := bs
	for {
		idx := bytes.Index(bs, key)
		if idx == -1 {
			break
		}
		bs = bs[idx+keylen:]
		idx = bytes.IndexFunc(bs, func(r rune) bool {
			return !unicode.IsNumber(r) && !unicode.IsLetter(r) && r != '_'
		})
		proto := bs[:idx]
		m[string(proto)] = 1
	}

	for {
		idx := bytes.Index(suffixbs, suffixkey)
		if idx == -1 {
			break
		}
		tmp := suffixbs[:idx]
		tmpidx := bytes.LastIndexFunc(tmp, func(r rune) bool {
			return !unicode.IsNumber(r) && !unicode.IsLetter(r) && r != '_'
		})
		name := tmp[tmpidx+1:]
		m[string(name)] = 1
		suffixbs = suffixbs[idx+suffixkeylen:]
	}
}
