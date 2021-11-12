package sys

import (
	"path/filepath"
	"sync"

	"k8s.io/klog/v2"
)

//LoaderType 加载器加载类型
const (
	LTFolder   = iota //NOTE: 加载类型-文件夹(遍历并根据文件名后缀调用对应处理，遍历子目录)
	LTFileYaml        //NOTE: 加载类型-yaml文件(不判断文件名后缀)
	LTFileJson        //NOTE: 加载类型-json文件(不判断文件名后缀)
)

type Config interface{}

var defaultLoadMux sync.Map

//HandleFolder register given Configs to given director (contains child directors)
func HandleFolder(dir string, v ...Config) {
	defer recoverWarning("sys: HandleFolder")
	handle(LTFolder, dir, v...)
}

//HandleJson register given Configs to given json file path
//NOTE: user should make sure the file of path is json type
func HandleJson(path string, v ...Config) {
	defer recoverWarning("sys: HandleJson:")
	handle(LTFileJson, path, v...)
}

//HandleYaml register given Configs to given yaml file path
//NOTE: user shoud make sure the file of path is yaml type
func HandleYaml(path string, v ...Config) {
	defer recoverWarning("sys: HandleYaml")
	handle(LTFileYaml, path, v...)
}

func recoverWarning(keymsg string) {
	if r := recover(); r != nil {
		klog.Warning(keymsg, r.(error))
	}
}

func handle(lt int, path string, v ...Config) {
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	var loader *Loader
	if val, ok := defaultLoadMux.Load(path); ok {
		loader = val.(*Loader)
	} else {
		loader = &Loader{lt: lt, path: path, vals: make([]Config, 0)}
		defaultLoadMux.Store(path, loader)
	}
	loader.Handle(v...)
}

func Unmarshal() {
	defaultLoadMux.Range(func(key, value interface{}) bool {
		loader := value.(*Loader)
		loader.Load()
		return true
	})
}
