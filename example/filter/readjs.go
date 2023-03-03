package filter

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

const (
	tagInote       = 1      //处于注释中
	tagInstring    = 1 << 1 //处于文字中
	tagInFat       = 1 << 2 //可换行
	tagInoteFat    = tagInFat | tagInote
	tagInstringFat = tagInFat | tagInstring
	// bitagIgnore       = 1
	// bitagIgnoreNote   = (1 << 1) | bitagIgnore //注释
	// bitagIgnoreString = (1 << 2) | bitagIgnore
	// bitagIgnoreLine   = (1 << 3) | bitagIgnoreNote
	// bitagIgnorePart   = (1 << 4) | bitagIgnoreNote
)

// func ReadJs(path string) []byte {
// 	bs, _ := os.ReadFile(path)
// 	// io.ByteReader
// 	reader := bufio.NewReader(bytes.NewBuffer(bs))
// 	bs = make([]byte, 0, len(bs))
// 	var ignore int
// L:
// 	for {
// 		tmp, _, err := reader.ReadLine()
// 		// if !isPrefix {
// 		// 	fmt.Println(tmp)
// 		// }
// 		switch err {
// 		case nil:
// 		case io.EOF:
// 			break L
// 		default:
// 			panic(err)
// 		}
// 		tmp = bytes.TrimSpace(tmp)
// 		if ignore & bitagIgnore != 0 {

// 			ignore = !bytes.Contains(tmp, []byte("*/"))
// 			continue
// 		}
// 		idx := bytes.IndexFunc(tmp, func(r rune) bool {
// 			switch rune {
// 			case "//":
// 				ignore = bitagIgnoreLine
// 			case "/*":
// 				ignore = bitagIgnorePart
// 			case '"':
// 				ignore = bitagIgnoreString
// 			}
// 			return ignore != 0
// 		})
// 		bs = append(bs, tmp[:idx]...)
// 		if ignore == bitagIgnoreString {
// 			bs = append(bs, tmp[0])
// 			tmp = tmp[1:]
// 			for {
// 				idx := bytes.IndexByte(tmp, '"')
// 				if idx == -1 {
// 					bs = append(bs, tmp...)
// 					continue L
// 				}
// 				if idx == 0 {

// 				}
// 			}
// 		}
// 		if idx := bytes.Index(tmp, []byte("//")); idx != -1 {
// 			bs = append(bs, bytes.TrimSpace(tmp[:idx])...)
// 			continue
// 		}
// 		if idx := bytes.Index(tmp, []byte("/*")); idx != -1 {
// 			ignore = !bytes.Contains(tmp[idx+2:], []byte("*/"))
// 			continue
// 		}
// 		bs = append(bs, tmp...)
// 	}
// 	return bs
// }

func ReadJs(path string) []byte {
	// fmt.Println(path)
	bs, _ := os.ReadFile(path)
	r := bufio.NewReader(bytes.NewReader(bs))
	rlt := make([]byte, 0, len(bs))
	var tag int // /**/ | ``
L:
	for {
		line, _, err := r.ReadLine()
		switch err {
		case nil:
		case io.EOF:
			break L
		default:
			panic(err)
		}

		for len(line) > 0 {
			// fmt.Println("--------------", tag)
			// fmt.Println(string(line))
			if tag&tagInFat == tagInFat { //可换行
				var valid []byte
				valid, line, tag = parseLineInFat(line, tag)
				if len(valid) > 0 {
					rlt = append(rlt, valid...)
				}
				if len(line) == 0 {
					continue L
				}
				////此时tag 一定是0，tag不为0的话，line一定为nil
			}
			var valid []byte
			valid, line, tag = parseLine(line)
			if len(valid) > 0 {
				rlt = append(rlt, valid...)
			}
		}
	}
	return rlt
}

func parseLine(line []byte) (valid, remain []byte, tag int) {
	idx := bytes.IndexFunc(line, func(r rune) bool {
		////无法解析代码中的 `` 符号，需要 \
		return r == '/' || r == '"' || r == '\''
	})
	if idx == -1 {
		valid = line
		return
	}
	switch line[idx] {
	case '/':
		if len(line) > idx+1 {
			valid = line[:idx]
			switch line[idx+1] {
			case '*':
				remain, tag = line[idx+2:], tagInoteFat
				return
			case '/':
				return
			}
		}
		valid = line
	case '"', '\'':
		sep := line[idx]
		valid, remain = line[:idx+1], line[idx+1:]
		for {
			ridx := bytes.IndexFunc(remain, func(r rune) bool {
				return r == '\\' || r == rune(sep)
			})
			if ridx == -1 {
				panic("invalid code")
			}
			if remain[ridx] == sep {
				valid, remain = append(valid, remain[:ridx+1]...), remain[ridx+1:]
				break
			}
			valid, remain = append(valid, remain[:ridx+2]...), remain[ridx+2:]
		}
	// case '\'':
	// 	valid, remain = line[:idx+3], line[idx+3:]
	// 	if line[idx+2] != '\'' {
	// 		fmt.Println("++++++++++++\n", string(valid))
	// 		panic("invalid code")
	// 	}
	default:
		panic("invalid bytes")
	}
	return
}

func parseLineInFat(line []byte, tag int) (valid, remain []byte, ntag int) {
	if tag&tagInFat == 0 {
		panic("not in fat")
	}
	if tag&tagInstring == tagInstring {
		idx := bytes.IndexByte(line, '`')
		if idx == -1 {
			valid, ntag = line, tag
			return
		}
		return line[:idx+1], line[idx+1:], 0
	}
	idx := bytes.Index(line, []byte("*/"))
	if idx == -1 {
		ntag = tag
		return
	}
	remain, ntag = line[idx+2:], 0
	return
}

// func splitRightMark(bs []byte) [][]byte {
// 	tmp := bs
// 	for {
// 		idx := bytes.IndexByte(tmp, '"')
// 		if idx == -1 {
// 			panic("error decode string")
// 		}
// 		tmp = tmp[idx+1:]
// 		if idx == 0 || tmp[idx-1] != '\\' {
// 			break
// 		}
// 	}
// 	return [][]byte{
// 		bs[:len(bs)-len(tmp)-1],
// 		tmp,
// 	}
// }
