package setting

import (
	"strings"
	"sync"
	"errors"
	"io/ioutil"
	"github.com/hiank/conf"
)


var _mtx sync.RWMutex

//LoadFromFile 从配置文件中读取需要的配置
func LoadFromFile(out interface{}, path string) error {

	_mtx.Lock()
	defer _mtx.Unlock()

	dotIdx := strings.LastIndexByte(path, '.')
	if dotIdx == -1 {
		return errors.New("file should be end with extension name")
	}
	extensionName := path[dotIdx+1:]		//NOTE: 扩展名

	var c conf.Conf
	switch strings.ToLower(extensionName) {
	case "json": c = conf.JSON
	case "yaml": c = conf.YAML
	default: return errors.New("not support file with extension name : " + extensionName)
	}

	in, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return c.Unmarshal(in, out)			//NOTE: 将文件内容解析到配置数据中
}