package cluster

import (
	"encoding/binary"
	"bytes"
)

// clusterType 
const (

	TypeKubIn 	= iota 	// kubernetes cluster in
	TypeKubOut			// kubernetes cluster out
) 

// GetAddr get server's addr[ip:port] by clusterType and msgName
func GetAddr(clusterType int, msgName string) (addr string, err error) {


	var ip string
	var port int32
	switch (clusterType) {

	case TypeKubIn: 
		ip, port, err = GetIPAndPortInKub(GetInClientset(), msgName)
	case TypeKubOut:
		// ip, port = GetIPAndPortInKub(GetOutClientset(), msgName)
	}

	if err == nil {
		buf := bytes.NewBuffer([]byte{})
		buf.WriteString(ip)
		binary.Write(buf, binary.BigEndian, port)
		addr = buf.String()	
	}
	return
}

