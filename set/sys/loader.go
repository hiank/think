package sys

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

var (
	JsonSuffix = "json"
	YamlSuffix = "yaml"
)

func suffix(path string) (val string) {
	if idx := strings.LastIndexByte(path, '.'); idx != -1 {
		val = strings.ToLower(path[idx+1:])
	}
	return
}

type Loader struct {
	path string
	lt   int      //NOTE: LoaderType
	vals []Config //NOTE: 需要从此配置中获取设置的所有配置对象
}

//Load load Configs
//NOTE: the func was synchronize
func (l *Loader) Load() {
	switch l.lt {
	case LTFileJson:
		l.loadJson(l.path)
	case LTFileYaml:
		l.loadYaml(l.path)
	case LTFolder:
		for _, path := range l.match() {
			switch suffix(path) {
			case JsonSuffix:
				l.loadJson(path)
			case YamlSuffix:
				l.loadYaml(path)
			default:
				klog.Warningf("sys: Load: cannot support file suffix with %s", suffix(path))
			}
		}
	default:
		klog.Warningf("sys: Load: unsupport LoaderType %d", l.lt)
	}
}

//Handle register given Configs
func (l *Loader) Handle(vals ...Config) {
	l.vals = append(l.vals, vals...)
}

//match 匹配所有文件(LTFolder LTFolderDep 用于获取所有需要读取的文件路径)
func (l *Loader) match() []string {
	dpaths, fpaths := l.listPaths(l.path)
	return l.tailMatch(dpaths, fpaths)
}

//tailMatch tial call algorithm for match all files
func (l *Loader) tailMatch(dpaths, fpaths []string) []string {
	ds, fs := l.listPaths(dpaths[0])
	dpaths = append(dpaths[1:], ds...)
	fpaths = append(fpaths, fs...)
	if len(dpaths) == 0 {
		return fpaths
	}
	return l.tailMatch(dpaths, fpaths)
}

//listPaths list folders and files in given folder path
//dpaths director paths
//fpaths file paths
func (l *Loader) listPaths(dpath string) (dpaths, fpaths []string) {
	fpaths, dpaths = make([]string, 0), make([]string, 0)
	dpath, err := filepath.Abs(dpath)
	if err != nil {
		klog.Warning("sys: Loader-listPaths: folder path invalid ", err)
		return
	}
	fis, err := ioutil.ReadDir(dpath)
	if err != nil {
		klog.Warning("sys: Loader-listPaths: ReadDir ", err)
		return
	}
	dir := dpath + string(os.PathSeparator)
	for _, fi := range fis {
		if fi.IsDir() {
			dpaths = append(dpaths, dir+fi.Name())
		} else {
			fpaths = append(fpaths, dir+fi.Name())
		}
	}
	return
}

//loadYaml read yaml file and set Configs
func (l *Loader) loadYaml(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		klog.Warning("sys: loadYaml: ", err)
		return
	}
	for _, cfg := range l.vals {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			klog.Warning("sys: loadYaml: yaml.Unmarshal to %v %v", cfg, err)
		}
	}
}

//loadJson read json file and set Configs
func (l *Loader) loadJson(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		klog.Warning("sys: loadJson: ", err, path)
		return
	}
	for _, cfg := range l.vals {
		if err := json.Unmarshal(data, cfg); err != nil {
			klog.Warningf("sys: loadJson: json.Unmarshal to %v %v", cfg, err)
		}
	}
}
