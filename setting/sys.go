package setting

import (
	"sync"
	"github.com/hiank/conf"
)


var defaultSysConf = `
{
	"sys.wsPort": 8022,
	"sys.k8sPort": 8026,
	"sys.maxMessageGo": 1000,
	"sys.clearInterval": 10
}
`

//Sys 系统配置
type Sys struct {

	WsPort 				int16 			`json:"sys.wsPort"`					//NOTE: websocket 服务端口
	K8sPort 			int16 			`json:"sys.k8sPort"`				//NOTE: k8s 服务端口
	MaxMessageGo 		int 			`json:"sys.maxMessageGo"`			//NOTE: MessageHub 的最大goruntine数量
	ClearInterval 		int 			`json:"sys.clearInterval"`			//NOTE: connhub 执行清理的时间间隔
}


var _sys *Sys
var _sysMtx sync.RWMutex
//GetSys 获得系统配置
func GetSys() *Sys {

	_sysMtx.Lock()
	defer _sysMtx.Unlock()

	if _sys == nil {
		_sys = new(Sys)
		c := conf.Conf(conf.JSON)
		c.Unmarshal([]byte(defaultSysConf), &_sys)
	}
	return _sys
}

