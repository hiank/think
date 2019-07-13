package conf

import (
	"strings"
	"reflect"
	// "path/filepath"
	// "os"
	"encoding/json"
	"io/ioutil"
	"sync"
)


//Item 配置管理类
type Item struct {

	Val 			interface{}
	reflectVal 		reflect.Value
}


//ValueByName 获得值
func (i *Item) ValueByName(filedName string) (r interface{}) {

	filedName = strings.ToUpper(filedName[:1]) + filedName[1:]

	t := i.reflectVal.Type()
	if _, ok := t.FieldByName(filedName); ok {

		r = i.reflectVal.FieldByName(filedName).Interface()
	}
	return 
}

//DoReflect 生成reflectVal 以便操作
func (i *Item) DoReflect() {

	i.reflectVal = reflect.ValueOf(i.Val)
	if i.reflectVal.Kind() == reflect.Ptr {
		i.reflectVal = i.reflectVal.Elem()
	}
}


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

	WsPort 				int64 			`json:"sys.wsPort"`					//NOTE: websocket 服务端口
	K8sPort 			int64 			`json:"sys.k8sPort"`				//NOTE: k8s 服务端口
	MaxMessageGo 		int 			`json:"sys.maxMessageGo"`			//NOTE: MessageHub 的最大goruntine数量
	ClearInterval 		int64 			`json:"sys.clearInterval"`			//NOTE: connhub 执行清理的时间间隔
}


// conf is json format setting
type conf struct {
	
	m 	map[string]*Item 		//NOTE: map[key]*Item
}

//Get 根据名称获得Item
func (c *conf) get(ns string) (i *Item, ok bool) {

	ns = strings.ToUpper(ns[:1]) + ns[1:]
	i, ok = c.m[ns]
	return
}

//Load 根据配置加载设置
func (c *conf) load(ci *Info) (err error) {
	
	v := ci.GetVal()
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {

		t = t.Elem()
	}
	ns := t.Name()

	i, ok := c.m[ns]
	if !ok {

		i = &Item{Val: v}
		c.m[ns] = i
	}
	
	buf := ci.GetBuffer()
	l := []*Item{i}
	if ns != "Sys" {

		if sys, ok := c.m["Sys"]; ok {

			l = append(l, sys)
		}
	}

	for _, i := range l {

		err = json.Unmarshal(buf, i.Val)
		i.DoReflect()	
	}
	return
}


//Info 配置Conf 文件
type Info struct {

	// path 		string
	buf 		[]byte
	val 		interface{}
}

//NewInfoByFile 新建一个ConfInfo
func NewInfoByFile(p string, v interface{}) (i *Info, e error) {

	var b []byte
	if b, e = ioutil.ReadFile(p); e == nil {

		i = &Info{b, v}
	}
	return
}

//NewInfo 创建配置信息
func NewInfo(b []byte, v interface{}) *Info {

	return &Info{b, v}
}

//GetBuffer 获得配置文件的绝对路径
func (ci *Info) GetBuffer() []byte {

	return ci.buf
}

//GetVal 获得配置文件的数据类型对象，不能是指针
func (ci *Info) GetVal() interface{} {

	return ci.val
}



var _conf *conf
var mtx sync.RWMutex
//Init 初始化k8s conn pool
func Init(cis ...*Info) {

	mtx.Lock()

	if _conf == nil {

		_conf = &conf{make(map[string]*Item)}
		_conf.load(NewInfo([]byte(defaultSysConf), &Sys{}))
	}
	for _, ci := range cis {

		_conf.load(ci)
	}

	mtx.Unlock()
}

// Get 获得全局Conf
func Get(ns string) *Item {

	mtx.RLock()

	item, _ := _conf.get(ns)
	
	mtx.RUnlock()
	return item
}

// GetSys 获得框架系统设置
func GetSys() *Sys {

	mtx.RLock()

	if _conf == nil {

		_conf = &conf{make(map[string]*Item)}
		_conf.load(NewInfo([]byte(defaultSysConf), &Sys{}))
	}
	item, _ := _conf.get("sys")

	mtx.RUnlock()
	return item.Val.(*Sys)
}

