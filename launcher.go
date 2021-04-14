package think

import (
	"context"
	"errors"
	"fmt"
	"sort"
)

var defaultLauncher = &launcher{}

type launcher struct {
	cache        []Module //NOTE: 缓存所有期望的Module
	createdCache []Module //NOTE: 缓存成功调用了OnCreate的Module
	startedCache []Module //NOTE: 缓存成功调用了OnStart的Module
}

func (le *launcher) create(ctx context.Context) error {
	for i, mod := range le.cache {
		if err := mod.OnCreate(ctx); err != nil {
			le.createdCache = le.cache[:i]
			return err
		}
	}
	le.createdCache = le.cache
	return nil
}

func (le *launcher) start(ctx context.Context) error {
	for i, mod := range le.cache {
		if err := mod.OnStart(ctx); err != nil {
			le.startedCache = le.cache[:i]
			return err
		}
	}
	le.startedCache = le.cache
	return nil
}

func (le *launcher) stop() {
	for _, mod := range le.startedCache {
		mod.OnStop()
	}
}

func (le *launcher) destroy() {
	for _, mod := range le.createdCache {
		mod.OnDestroy()
	}
}

type sortWeight []*weight

func (sw sortWeight) Len() int {
	return len(sw)
}

func (sw sortWeight) Less(i, j int) bool {
	return sw[i].toInt() > sw[j].toInt()
}

func (sw sortWeight) Swap(i, j int) {
	sw[i], sw[j] = sw[j], sw[i]
}

//Launch 处理模块启动
//每个进程只能调用一次
func Launch(mods ...Module) (err error) {
	if defaultLauncher.cache != nil {
		return errors.New("already launched")
	}
	weightMap := make(map[Module]*weight)
	for _, mod := range mods {
		tagMap := make(map[Module]bool)
		if err := calcWeight(weightMap, mod, tagMap); err != nil {
			return err
		}
	}

L:
	for key := range weightMap {
		for _, mod := range mods {
			if mod == key {
				continue L
			}
		}
		return fmt.Errorf("required %v not launch", key)
	}

	ctx := context.Background()
	defaultLauncher.cache = modulesSorted(weightMap)
	if err = defaultLauncher.create(ctx); err == nil {
		err = defaultLauncher.start(ctx)
	}
	return
}

func Unlaunch() {
	defaultLauncher.stop()
	defaultLauncher.destroy()
	defaultLauncher = &launcher{}
}

type weight struct {
	mod     Module
	value   int
	depends []*weight
}

func newWeight(mod Module, val int) *weight {
	w := &weight{
		mod:     mod,
		value:   val,
		depends: make([]*weight, 0, 8),
	}
	return w
}

func (w *weight) toInt() int {
	return w.value
}

func (w *weight) add(val int) {
	w.value += val
	for _, depend := range w.depends {
		depend.add(val)
	}
}

func (w *weight) addDependVal(depend *weight) {
	w.depends = append(w.depends, depend)
}

func calcWeight(weightMap map[Module]*weight, mod Module, tagMap map[Module]bool) error {
	lastWeight, ok := weightMap[mod]
	if !ok {
		lastWeight = newWeight(mod, -1)
		weightMap[mod] = lastWeight
	}

	tagMap[mod] = true //NOTE: 用于标识依赖包链
	for _, dependMod := range mod.Depend() {
		if tagMap[dependMod] {
			return fmt.Errorf("loop depend %v", dependMod)
		}
		weight, ok := weightMap[dependMod]
		if ok {
			if weight.toInt() <= lastWeight.toInt() {
				weight.add(lastWeight.toInt() - weight.toInt() + 1)
			}
			lastWeight.addDependVal(weight)
			continue
		}
		weight = newWeight(dependMod, lastWeight.toInt()+1)
		weightMap[dependMod] = weight
		lastWeight.addDependVal(weight)
		if err := calcWeight(weightMap, dependMod, newTagMap(tagMap)); err != nil {
			return err
		}
	}
	return nil
}

func newTagMap(tagMap map[Module]bool) map[Module]bool {
	newMap := make(map[Module]bool)
	for key, val := range tagMap {
		newMap[key] = val
	}
	return newMap
}

func modulesSorted(weightMap map[Module]*weight) []Module {
	idx, sw, mods := 0, make([]*weight, len(weightMap)), make([]Module, len(weightMap))
	for _, w := range weightMap {
		sw[idx] = w
		idx++
	}
	sort.Sort(sortWeight(sw))

	for idx, w := range sw {
		mods[idx] = w.mod
	}
	return mods
}
