package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	m := readXNpc()
	// fmt.Println(m)

	root := "path/"
	os.Mkdir(root, 0777)
	//
	for key, val := range m {
		//
		os.Mkdir(root + key, 0777)
		os.Mkdir(root + key + "/" + val, 0777)
	}
}

func readXNpc() map[string]string {
	csvFile, err := os.Open("./XNpc.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	r := csv.NewReader(bufio.NewReader(csvFile))
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	head := records[0]
	records = records[1:]
	fmt.Println(len(records))
	// codes, eftIds := make([]string, 0, len(records)), make([]string, 0, len(records))
	var codeIdx, eftIdsIdx, cnt int
	for idx, key := range head {
		//
		if key == "code" {
			codeIdx = idx
		} else if key == "eftId" {
			eftIdsIdx = idx
		} else {
			continue
		}
		cnt++
		if cnt == 2 {
			break
		}
	}
	m := make(map[string]string)
	for _, val := range records {
		// codes, eftIds = append(codes, val[codeIdx]), append(eftIds, val[eftIdsIdx])
		code := val[codeIdx]
		if _, ok := m[code]; ok {
			fmt.Println(code)
		}
		m[code] = val[eftIdsIdx]
	}
	fmt.Println(len(m))
	return m
}
