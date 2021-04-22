package modex

import (
	"errors"
	"strconv"
	"strings"
)

type Addr struct {
	Value string `json:"Addr"`
}

type AddrGetter interface {
	Get() *Addr
}

type AddrParser interface {
	AutoAddr(serveName, portName string) (string, error)
}

//ParseAddr 从配置的Addr中解析包含端口的最终Addr
func ParseAddr(addr string, parser AddrParser) (url string, err error) {
	if addr == "" {
		return "", errors.New("addr should not be empty")
	}

	arr := strings.Split(addr, ":")
	switch len(arr) {
	case 1:
		url, err = parser.AutoAddr(arr[0], "")
	case 2:
		url = addr
		if _, numErr := strconv.Atoi(arr[0]); numErr != nil {
			url, err = parser.AutoAddr(arr[0], arr[1])
		}
	default:
		err = errors.New("too many sep for Addr: " + addr)
	}
	return
}
