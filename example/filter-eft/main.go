package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hiank/think/doc"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func main() {

	bts, _ := os.ReadFile("./config.json")
	// coder := doc.NewCoder[doc.JsonCoder]()

	var coder doc.Json = bts
	var cfg config
	coder.Decode(&cfg)

	m := make(map[string]string)
	for _, folder := range cfg.EftFolders {
		readEfts(folder, m)
	}
	////
	keys := maps.Keys(m)
	slices.Sort(keys)
	datafull := make([]byte, 0, 1024)
	for _, path := range keys {
		datafull = append(datafull, path...)
		datafull = append(datafull, '\n')
	}
	fmt.Println(".+++++++++++++++", len(keys), "----", len(m))
	os.WriteFile("./allpaths.txt", datafull, 0777)
	// return

	src := make(map[string][]byte)
	for _, folder := range cfg.JsFolders {
		readCodes(folder, cfg.JsIgnores, src)
	}
	for _, folder := range cfg.CsvFolders {
		readCsvs(folder, src)
	}

	data, data2 := make([]byte, 0, 1024), make([]byte, 0, 1024)
L:
	for path, name := range m {
		idx := strings.LastIndexByte(name, '.')
		tmpName := name[:idx]
		// fmt.Println(tmpName)
		for jspath, b := range src {
			idx := strings.Index(string(b), tmpName)
			if idx != -1 {
				fmt.Printf("%s in %s\n", name, jspath)
				continue L
			}
		}

		data = append(data, path...)
		data = append(data, '\n')

		data2 = append(data2, name...)
		data2 = append(data2, '\n')
	}

	os.WriteFile("pathes.txt", data, 0777)
	os.WriteFile("names.txt", data2, 0777)
}

func readCsvs(folder string, out map[string][]byte) {
	filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			out[path], _ = os.ReadFile(path)
		}
		return nil
	})
}

func readCodes(folder string, ignores []string, out map[string][]byte) {
	filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".js" {
			// for _, ignore := range ignores {
			////
			// n3, _ := filepath.Abs(ignore)
			// if n2 == n3 {
			// 	return nil
			// }
			// }
			out[path], _ = os.ReadFile(path)
		}
		return nil
	})
}

func readEfts(folder string, out map[string]string) {
	filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		// fmt.Println(path)
		n1, _ := filepath.Abs(folder)
		n2, _ := filepath.Abs(path)
		if n1 == n2 {
			return nil
		}
		if info.IsDir() {
			readEfts(path, out)
		} else if suffixWithEft(path) {
			out[path] = info.Name()
		}
		return nil
	})
}

func withSuffix(path, suffix string) (suc bool) {
	idx := strings.LastIndexByte(path, '.')
	if idx != -1 {
		suc = path[idx+1:] == suffix
	}
	return
}

func suffixWithEft(path string) (suc bool) {
	return withSuffix(path, "atlas") || withSuffix(path, "ExportJson")
}

func suffixWithJs(path string) (suc bool) {
	return withSuffix(path, "js")
}

type config struct {
	EftFolders []string `json:"folders.eft"`
	JsFolders  []string `json:"folders.js"`
	CsvFolders []string `json:"folders.csv"`
	JsIgnores  []string `json:"ignores.js"`
}
