package filter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// type csvJson struct {
// 	M map[string][]string `json:""`
// }

// func ScanCsvInFolder()

func ScanCsvJsAndFiles(jsPath, csvCopyFolder string, csvFileM map[string]string) (filenames, texts []string) {
	csvFileM = maps.Clone(csvFileM)
	bs, _ := os.ReadFile(jsPath)
	idx := bytes.Index(bs, []byte("Config = "))
	if idx == -1 {
		panic("invalid js file")
	}
	bs = bs[idx:]
	idx = bytes.IndexByte(bs, '{')
	bs = bs[idx+1:]
	idx = bytes.IndexByte(bs, '}')
	bs = bs[:idx]

	bio := bufio.NewReader(bytes.NewBuffer(bs))
	// r := bytes.NewBuffer(bs)
	injs := make(map[string]int)
L:
	for {
		line, _, err := bio.ReadLine()
		switch err {
		case nil:
		case io.EOF:
			break L
		default:
			panic("invalid js")
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] == '/' {
			continue L
		}
		arr := bytes.Split(line, []byte{':'})
		switch len(arr) {
		case 1:
			continue L
		case 2:
		default:
			fmt.Println("-----------\n", string(line))
			panic("invalid line")
		}
		name := string(bytes.TrimSpace(arr[0]))
		fmt.Printf(":%s\n", name)
		injs[name] = 1
	}
	// scanCsvInFolder(csvFolder, mfile)

	keys := maps.Keys(injs)
	for _, key := range keys {
		if path, ok := csvFileM[key]; ok {
			delete(injs, key)
			delete(csvFileM, key)

			// exec.Command("cp", path, filepath.Join(csvCopyFolder, key+".csv"))
			copyFile(filepath.Join(csvCopyFolder, key+".csv"), path)
		}
	}

	filenames, texts = maps.Keys(csvFileM), maps.Keys(injs)
	slices.Sort(filenames)
	slices.Sort(texts)
	return
}

func copyFile(dstPath, srcPath string) {
	src, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		panic(err)
	}
	defer dst.Close()
	io.Copy(dst, src)
}

func ScanCsvInFolder(folder string, m map[string]string) {
	filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			if name := strings.TrimSuffix(info.Name(), ".csv"); len(name) != len(info.Name()) {
				m[name] = path
			}
		}
		return nil
	})
}

func ScanNotUsedInJs(jsFolder string, m map[string]string, jsIgnores []string) (notused map[string]string) {
	// notused = make(map[string]string)
	m = maps.Clone(m)
	// sm := StringMap{
	// 	Spine:   make(map[string]int),
	// 	Export:  make(map[string]int),
	// 	Special: make(map[string]int),
	// }
	// readJstringsInFolder(jsFolder, jsIgnores, sm)
	for i, path := range jsIgnores {
		jsIgnores[i], _ = filepath.Abs(path)
	}
	slices.Compact(jsIgnores)
	WalkGivenExts(jsFolder, func(path string) error {
		ps, _ := filepath.Abs(path)
		if !slices.Contains(jsIgnores, ps) {
			scanCsvInJs(path, m)
		}
		return nil
	}, ".js")
	return m
}

func scanCsvInJs(jsPath string, m map[string]string) {
	bs, keys := ReadJs(jsPath), maps.Keys(m)
	prefix := []byte("Config.")
	keys = slices.Compact(keys)
	slices.SortFunc(keys, func(a, b string) bool {
		return b < a
	})
	for _, key := range keys {
		tmpBs := bs
		kb := append(prefix, []byte(key)...)
		idx := bytes.Index(tmpBs, kb)
	L2:
		for idx != -1 {
			if len(tmpBs) > idx+len(kb) {
				r := rune(tmpBs[idx+len(kb)])
				if unicode.IsDigit(r) || unicode.IsLower(r) || unicode.IsUpper(r) {
					tmpBs = tmpBs[idx+len(kb):]
					idx = bytes.Index(tmpBs, kb)
					continue L2
				}
			}
			// out[key] = m[key]
			delete(m, key)
			break L2
		}
	}
}
