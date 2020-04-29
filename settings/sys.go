package settings

import (
	"sync"

	"github.com/hiank/conf"
)

var defaultSysConf = `
{
	"sys.wsPort": 8022,
	"sys.k8sPort": 8026,
	"sys.messageGo": 1000,
	"sys.grpcGo": 20,
	"sys.clearInterval": 10,
	"sys.timeOut": 360
}
`

//Sys 系统配置
type Sys struct {
	WsPort        uint16 `json:"sys.wsPort"`        //NOTE: websocket 服务端口
	K8sPort       uint16 `json:"sys.k8sPort"`       //NOTE: k8s 服务端口
	MessageGo     int    `json:"sys.messageGo"`     //NOTE: MessageHub 的最大goroutine数量
	GrpcGo        int    `json:"sys.grpcGo"`        //NOTE: grpc 最大连接数，每个连接对应一个goroutine
	ClearInterval int    `json:"sys.clearInterval"` //NOTE: connhub 执行清理的时间间隔
	TimeOut       int64  `json:"sys.timeOut"`       //NOTE: 连接无响应最大时间, 单位s
}

var _sys *Sys
var _sysOnce sync.Once

//GetSys 获得系统配置
func GetSys() *Sys {

	_sysOnce.Do(func() {
		_sys = new(Sys)
		c := conf.Conf(conf.JSON)
		c.Unmarshal([]byte(defaultSysConf), &_sys)
	})
	return _sys
}
