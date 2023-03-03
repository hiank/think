package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

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

	m := make(map[string]int)
	for _, name := range cfg.Paths {
		filterImages(cfg.Root+name, m)
	}

	for _, folder := range cfg.Folders {
		fmt.Println(folder)
		filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				fmt.Println(info.Name())
				filterImages(path, m)
			}
			return nil
		})
	}

	// io.
	// ioutil.
	// for path, _ := range m {

	// }
	keys := maps.Keys(m)
	slices.Sort(keys)
	data := make([]byte, 0, 1024)
	for _, path := range keys {
		data = append(data, path...)
		data = append(data, '\n')
	}

	docpath := cfg.Docment
	if docpath == "" {
		docpath = "./tmp.txt"
	}

	os.WriteFile(docpath, data, 0777)

	fmt.Printf("%v\n", docpath)
}

func filterImages(path string, out map[string]int) {
	///
	bts, _ := os.ReadFile(path)
	for idx, key := 0, []byte("Path="); idx != -1; idx = bytes.Index(bts, key) {
		bts = bts[idx+len(key)+1:]
		endIdx := bytes.IndexByte(bts, '"')
		if v := bts[:endIdx]; len(v) > 4 {
			if tail := string(v[len(v)-4:]); tail == ".png" || tail == ".jpg" {
				out[string(v)] = 0
			}
		}
		bts = bts[endIdx+1:]
	}
}

type config struct {
	Root    string   `json:"root"`
	Paths   []string `json:"paths"`
	Docment string   `json:"docment"`
	Folders []string `json:"folders"`
}
