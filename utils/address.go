package utils

import (
	"bytes"
	"strconv"
)

//WithPort ip + port to string
func WithPort(ip string, port uint16) string {

	var buffer bytes.Buffer
	buffer.WriteString(ip)
	buffer.WriteByte(':')
	buffer.WriteString(strconv.FormatInt(int64(port), 10))
	return buffer.String()
}
