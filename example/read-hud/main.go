package main

import (
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

	m := make(map[string]string)
	// for _, name := range cfg.Paths {
	// 	filterImages(cfg.Root+name, m)
	// }

	wantFolder, _ := filepath.Abs(cfg.WantFolder)
	wfLen := len(wantFolder)
	filepath.Walk(wantFolder, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			// fmt.Println(info.Name())
			// filterImages(path, m)
			// os.Remove()
			// strings.In
			key := path[wfLen:]
			m[key] = path
		}
		return nil
	})

	// paths := make([]string, 0, 16)
	folders := make([]string, 0, 128)
	sFolder, _ := filepath.Abs(cfg.SFolder)
	sfLen := len(sFolder)
	filepath.Walk(sFolder, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			key := path[sfLen:]
			//
			if _, ok := m[key]; !ok {
				if err := os.Remove(path); err != nil {
					fmt.Println(err)
				}
			} else {
				delete(m, key)
			}
		} else {
			folders = append(folders, path)
		}
		return nil
	})

	slices.SortFunc(folders, func(p1, p2 string) bool {
		return p1 > p2
	})

	for _, path := range folders {
		fis, _ := os.ReadDir(path)
		if len(fis) == 0 { //空目录
			os.Remove(path)
		}
	}

	paths := maps.Values(m)
	slices.Sort(paths)

	data := make([]byte, 0, 1024)
	for _, path := range paths {
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

type config struct {
	Root       string   `json:"root"`
	Paths      []string `json:"paths"`
	Docment    string   `json:"docment"`
	SFolder    string   `json:"s-folder"`
	WantFolder string   `json:"want-folder"`
}
