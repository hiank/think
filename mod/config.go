package mod

import (
	"container/list"
	"context"
	"errors"

	"github.com/hiank/think"
	"github.com/hiank/think/set"
	"k8s.io/klog/v2"
)

var Config = &config{
	values:  make(map[string]*list.List),
	folders: []string{},
}

type config struct {
	folders []string              //NOTE: 配置文件目录组
	values  map[string]*list.List //NOTE: 期望从配置文件中获取数据的对象

	think.IgnoreDepend
	think.IgnoreOnCreate
	think.IgnoreOnStop
	think.IgnoreOnDestroy
}

func (ms *config) SignUpValue(key string, vals ...interface{}) {
	cache, ok := ms.values[key]
	if !ok {
		cache = list.New()
		ms.values[key] = cache
	}
	for _, val := range vals {
		if val != nil {
			cache.PushBack(val)
		}
	}
}

func (ms *config) SignUpFolder(folders ...string) (outErr error) {
	ms.folders = set.LookRootFolders(append(folders, ms.folders...))
	if len(ms.folders) == 0 {
		outErr = errors.New("no folder available")
	}
	return
}

func (ms *config) OnStart(ctx context.Context) error {
	if len(ms.folders) == 0 {
		return errors.New("config folder not set")
	}
	cacheMap := make(map[string]*list.List)
	for _, folder := range ms.folders {
		set.WalkText(cacheMap, folder, set.JSON, set.YAML)
	}
L:
	for key, cache := range cacheMap {
		valist, ok := ms.values[key]
		if !ok {
			continue L
		}
		switch key {
		case set.JSON:
			set.UnmarshalJSON(cache, valist)
		case set.YAML:
			fallthrough
		default:
			klog.Warningf("not yet support %v config\n", key)
		}
	}
	return nil
}
