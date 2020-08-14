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
	"sys.timeOut": 360000,
	"sys.TokDonedLen": 10,
	"sys.TokReqLen": 64,
	"sys.ConnHubReqLen": 64,
	"sys.MessageHubReqLen": 0
}
`

//Sys 系统配置
type Sys struct {
	WsPort           uint16 `json:"sys.wsPort"`           //NOTE: websocket 服务端口
	K8sPort          uint16 `json:"sys.k8sPort"`          //NOTE: k8s 服务端口
	MessageGo        int    `json:"sys.messageGo"`        //NOTE: MessageHub 的最大goroutine数量
	GrpcGo           int    `json:"sys.grpcGo"`           //NOTE: grpc 最大连接数，每个连接对应一个goroutine
	TimeOut          int64  `json:"sys.timeOut"`          //NOTE: 连接无响应最大时间, 单位ms
	TokDonedLen      int    `json:"sys.TokDonedLen"`      //NOTE: token失效清理chan 的缓存长度，避免清理请求阻塞
	TokReqLen        int    `json:"sys.TokReqLen"`        //NOTE: token请求chan 的缓存长度
	ConnHubReqLen    int    `json:"sys.ConnHubReqLen"`    //NOTE: connHub请求chan 的缓存长度
	MessageHubReqLen int    `json:"sys.MessageHubReqLen"` //NOTE: MessageHub请求chan 的缓存长度
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
