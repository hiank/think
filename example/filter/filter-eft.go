package filter

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type EftMap struct {
	Spine   map[string]string
	Export  map[string]string
	Special map[string]string
}

type StringMap struct {
	Spine   map[string]int
	Export  map[string]int
	Special map[string]int
}

// UnusedEfts map[path]name
func UnusedEfts(eftFolder, jsFolder, csvFolder string, jsIgnores, eftFolderIgnores []string) (ue map[string]string, strs []string) {
	em := EftMap{
		Spine:   make(map[string]string),
		Export:  make(map[string]string),
		Special: make(map[string]string),
	}
	sm := StringMap{
		Spine:   make(map[string]int),
		Export:  make(map[string]int),
		Special: make(map[string]int),
	}
	readJstringsInFolder(jsFolder, jsIgnores, sm)
	fmt.Printf("++++++++%v\n", sm.Spine["spine_zhandou_zjdz"])

	readCsvFolder(csvFolder, sm)

	readEftsInFolder(eftFolder, eftFolderIgnores, em)
	checkRune(em.Export)
	checkRune(em.Spine)
	checkRune(em.Special)
	// return map[string]string{}, []string{}
	m := make(map[string]string)
	scanEfts(em.Special, sm.Special, m)
	scanEfts(em.Export, sm.Export, m)
	scanEfts(em.Spine, sm.Spine, m)
	///
	names := maps.Keys(sm.Spine)
	names = append(names, maps.Keys(sm.Export)...)
	return m, names
}

func matchEftnameByte(name string) (match bool) {
	if strings.IndexFunc(name, func(r rune) bool {
		switch {
		case r == '_':
		case r >= '0' && r <= '9':
		case r >= 'A' && r <= 'Z':
		case r >= 'a' && r <= 'z':
		default:
			return true
		}
		return false
	}) == -1 {
		match = true
	}
	return
}

func checkRune(m map[string]string) {
	for _, name := range m {
		if !matchEftnameByte(name) {
			fmt.Println(name)
			panic("unsupport name")
		}
	}
}

func scanEfts(eftm map[string]string, strm map[string]int, unused map[string]string) {
	fmt.Printf(".......%d.....%d\n", len(eftm), len(strm))
	for path, name := range eftm {
		if _, ok := strm[name]; !ok {
			fmt.Println(name)
			unused[path] = name
		} else {
			fmt.Println("-=-=-=-=-=-=-=-=-=-=-=")
			delete(strm, name)
		}
	}
}

func readJstringsInFolder(folder string, ignores []string, sm StringMap) {
	// ignores = append(ignores, folder)
	for i, path := range ignores {
		ignores[i], _ = filepath.Abs(path)
	}
	slices.Compact(ignores)
	WalkGivenExts(folder, func(path string) error {
		ps, _ := filepath.Abs(path)
		if !slices.Contains(ignores, ps) {
			readJstringsInFile(path, sm)
		}
		return nil
	}, ".js")
}

func readJstringsInFile(path string, sm StringMap) {
	bs := ReadJs(path)
	for {
		idxl := bytes.IndexByte(bs, '"')
		if idxl == -1 {
			break
		}
		bs = bs[idxl+1:]
		idxr := bytes.IndexByte(bs, '"')
		if idxr == -1 {
			fmt.Println(path)
			fmt.Println(string(bs))
			panic("error js file")
		}
		unmarshalToStringMap(sm, string(bs[:idxr]))
		bs = bs[idxr+1:]
	}
}

func unmarshalToStringMap(sm StringMap, str string) {
	for {
		idx := strings.LastIndexByte(str, '.')
		if idx == -1 {
			break
		}
		if len(str) > idx+1 && str[idx+1] == '/' {
			str = str[idx+2:]
			break
		}
		str = str[:idx]
	}
	arr := strings.Split(str, "/")
	for _, str := range arr {
		if !matchEftnameByte(str) {
			continue
		}
		preidx := strings.IndexByte(str, '_')
		switch {
		case preidx == -1:
			sm.Special[str] = 1
		case str[:preidx] == "eft":
			sm.Export[str] = 1
		case str[:preidx] == "spine":
			sm.Spine[str] = 1
		default:
			// fmt.Println(str)
			sm.Special[str] = 1
		}
	}
}

func readCsvFolder(folder string, sm StringMap) {

	filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			readCsvFile(path, sm)
		}
		return nil
	})
}

func readCsvFile(path string, sm StringMap) {
	bs, _ := os.ReadFile(path)
	bytes.TrimFunc(bs, func(r rune) bool {
		return r == '\n' || r == ' '
	})
	bss := bytes.Split(bs, []byte{'\n'})
	if len(bss) < 2 {
		return
	}
	title := bss[0]
	bytes.TrimFunc(title, func(r rune) bool {
		return r == '\n' || r == ' ' || r == '\r'
	})
	num := len(bytes.Split(title, []byte{','}))
	for _, row := range bss[1:] {
		row = bytes.TrimFunc(row, func(r rune) bool {
			return r == '\n' || r == ' ' || r == '\r'
		})
		vals := bytes.Split(row, []byte{','})
		if len(vals) != num {
			fmt.Println(path)
			fmt.Println(string(title))
			fmt.Println(string(row))
			fmt.Printf("%v____%v\n", len(vals), num)
			// panic(".......")
		}
		for _, val := range vals[1:] {
			if len(val) > 0 {
				unmarshalToStringMap(sm, string(val))
			}
		}
	}
}

func hasPrefixDir(rootDir, curDir string) (suc bool) {
	rootDir, _ = filepath.Abs(rootDir)
	curDir, _ = filepath.Abs(curDir)
	curLen, rootLen := len(curDir), len(rootDir)
	if rootLen > curLen || curDir[:rootLen] != rootDir {
		return false
	}
	return (rootLen == curLen) || (curDir[rootLen] == filepath.Separator)
}

func readEftsInFolder(folder string, eftFolderIgnores []string, em EftMap) {
	// ignores = append(ignores, folder)
	// for i, path := range eftFolderIgnores {
	// 	eftFolderIgnores[i], _ = filepath.Abs(path)
	// }
	// slices.Compact(eftFolderIgnores)

	filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			for _, ignore := range eftFolderIgnores {
				if hasPrefixDir(ignore, filepath.Dir(path)) {
					return nil
				}
			}
			ext := filepath.Ext(path)
			switch strings.ToLower(ext) {
			case ".exportjson", ".atlas":
				name := info.Name()[:len(info.Name())-len(ext)]
				preidx := strings.IndexByte(name, '_')
				switch {
				case preidx == -1:
					em.Special[path] = name
				case name[:preidx] == "eft":
					em.Export[path] = name
				case name[:preidx] == "spine":
					em.Spine[path] = name
				default:
					em.Special[path] = name
				}
			}
		}
		return nil
	})
}
