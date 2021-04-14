package think_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hiank/think"
	"gotest.tools/v3/assert"
)

type testMod1 struct {
	key         int
	createdChan chan int
	depends     []think.Module
	Created     bool
	Started     bool
	Stopped     bool
	Destroyed   bool
	lockCreate  bool
	lockStart   bool
}

func (tm *testMod1) Depend() []think.Module {
	return tm.depends
}

func (tm *testMod1) OnCreate(context.Context) (err error) {
	tm.Created = true
	if tm.lockCreate {
		err = errors.New("create locked")
	}
	select {
	case tm.createdChan <- tm.key:
	default:
	}
	return err
}

func (tm *testMod1) OnStart(context.Context) (err error) {
	tm.Started = true
	if tm.lockStart {
		err = errors.New("start locked")
	}
	return err
}

func (tm *testMod1) OnStop() {
	tm.Stopped = true
}

func (tm *testMod1) OnDestroy() {
	tm.Destroyed = true
}

func TestLaunchSimple(t *testing.T) {

	tm := &testMod1{}
	err := think.Launch(tm)
	assert.Assert(t, err == nil, "")

	assert.Assert(t, tm.Created && tm.Started)
	assert.Equal(t, tm.Destroyed, false)

	think.Unlaunch()
	assert.Assert(t, tm.Destroyed)
}

func TestLaunchMoreTimes(t *testing.T) {

	tm := &testMod1{}
	err := think.Launch(tm)
	assert.Assert(t, err == nil)
	// think.Unlaunch()

	err = think.Launch(tm)
	assert.Assert(t, err != nil)

	think.Unlaunch()
}

func TestLaunchDependModules(t *testing.T) {

	tm1 := &testMod1{}
	tm2 := &testMod1{depends: []think.Module{tm1}}

	err := think.Launch(tm2)
	assert.Assert(t, err != nil, err)
	think.Unlaunch()

	err = think.Launch(tm1, tm2)
	assert.Assert(t, err == nil)
	think.Unlaunch()
}

//依赖顺序，被依赖的模块必须先被加载
func TestLaunchDependModulesOrder(t *testing.T) {

	created := make(chan int, 2)

	tm1 := &testMod1{key: 1, createdChan: created}
	tm2 := &testMod1{key: 2, createdChan: created, depends: []think.Module{tm1}}

	err := think.Launch(tm2, tm1)
	assert.Assert(t, err == nil, "加载会自动调整依赖顺序")

	assert.Equal(t, <-created, 1, "被依赖的Module需要先被加载")
	assert.Equal(t, <-created, 2)

	think.Unlaunch()
}

func TestLaunchDependModulesMissing(t *testing.T) {

	tm1 := &testMod1{}
	tm2 := &testMod1{depends: []think.Module{tm1}}

	err := think.Launch(tm2)
	assert.Assert(t, err != nil, "如果依赖的Module不在加载列表中，则加载失败")

	think.Unlaunch()
}

func TestLaunchCreateError(t *testing.T) {

	tm1, tm2 := &testMod1{lockCreate: true}, &testMod1{}
	err := think.Launch(tm1, tm2)
	assert.Assert(t, err != nil)

	assert.Assert(t, !tm1.Started)
	assert.Assert(t, !tm2.Started)

	think.Unlaunch()
	assert.Assert(t, !tm1.Stopped)
	assert.Assert(t, !tm2.Stopped)

	assert.Assert(t, !tm1.Destroyed)
	assert.Assert(t, !tm2.Destroyed)
}

func TestLaunchStartError(t *testing.T) {
	tm1, tm2 := &testMod1{}, &testMod1{lockStart: true}
	err := think.Launch(tm1, tm2)
	assert.Assert(t, err != nil)

	assert.Assert(t, tm1.Created)
	assert.Assert(t, tm2.Created)

	think.Unlaunch()
	assert.Assert(t, tm1.Destroyed)
	assert.Assert(t, tm2.Destroyed)

	assert.Assert(t, tm1.Stopped, "tm1是成功加载的，所以需要有个OnStop过程")
	assert.Assert(t, !tm2.Stopped, "tm2的OnCreate会返回错误，因此相应的OnStop过程不会被调用")
}

func TestMapNon(t *testing.T) {
	tmpMap := make(map[string]int)
	assert.Equal(t, tmpMap["ok"], 0, "未设值的字段，返回默认值")
}

func TestArrayEmpty(t *testing.T) {
	var data []int
	assert.Equal(t, len(data), 0)
}

func TestSortWeight(t *testing.T) {
	tms := make([]*testMod1, 20)
	for i := range tms {
		tms[i] = &testMod1{key: i}
	}

	// tms[11].depends = []Module{tms[3]}
	tms[3].depends = []think.Module{tms[1], tms[4], tms[7]}
	tms[4].depends = []think.Module{tms[7], tms[8]}
	tms[7].depends = []think.Module{tms[2]}
	tms[1].depends = []think.Module{tms[9]}
	tms[9].depends = []think.Module{tms[7], tms[10]}
	tms[10].depends = []think.Module{tms[7]}

	weightMap := think.Export_makeWeightMap()
	for _, tim := range tms {
		tagMap := make(map[think.Module]bool)
		err := think.Export_calcWeight(weightMap, tim, tagMap)
		assert.Assert(t, err == nil)
	}

	mods := think.Export_modulesSorted(weightMap)
	assert.Equal(t, len(mods), len(tms))

	var lastMod think.Module
	for _, mod := range mods {
		if lastMod == nil {
			lastMod = mod
			continue
		}
		assert.Assert(t, think.Export_weightValue(weightMap[lastMod]) >= think.Export_weightValue(weightMap[mod]), fmt.Sprintf("%v__%v", think.Export_weightValue(weightMap[lastMod]), think.Export_weightValue(weightMap[mod])))
		lastMod = mod
	}
}

func TestCalcWeight(t *testing.T) {
	tms := make([]*testMod1, 20)
	for i := range tms {
		tms[i] = &testMod1{key: i}
	}

	tms[3].depends = []think.Module{tms[1], tms[4], tms[7]}
	tms[4].depends = []think.Module{tms[7], tms[8]}
	tms[7].depends = []think.Module{tms[2]}
	tms[1].depends = []think.Module{tms[9]}
	tms[9].depends = []think.Module{tms[7], tms[10]}
	tms[10].depends = []think.Module{tms[7]}

	weightMap := think.Export_makeWeightMap()
	for _, tim := range tms {
		tagMap := make(map[think.Module]bool)
		err := think.Export_calcWeight(weightMap, tim, tagMap)
		assert.Assert(t, err == nil)
	}

	arr := [][]int{
		{1, 3},
		{4, 3},
		{7, 3},
		{7, 4},
		{8, 4},
		{2, 7},
		{9, 1},
		{7, 9},
		{10, 9},
		{7, 10},
	}

	for _, pair := range arr {
		w1, w2 := weightMap[tms[pair[0]]], weightMap[tms[pair[1]]]
		assert.Assert(t, think.Export_weightValue(w1) > think.Export_weightValue(w2), fmt.Sprintf("%v:%v__%v:%v", pair[0], think.Export_weightValue(w1), pair[1], think.Export_weightValue(w2)))
	}

	tms[11].depends = []think.Module{tms[3]}
	tms[8].depends = []think.Module{tms[3]}
	weightMap = think.Export_makeWeightMap()
	tagMap := make(map[think.Module]bool)
	err := think.Export_calcWeight(weightMap, tms[3], tagMap)
	assert.Assert(t, err != nil, err)

	tms[11].depends = []think.Module{}
	tms[8].depends = []think.Module{}
	weightMap = think.Export_makeWeightMap()
	tagMap = make(map[think.Module]bool)
	err = think.Export_calcWeight(weightMap, tms[11], tagMap)
	assert.Assert(t, err == nil, err)

	tms[11].depends = []think.Module{&testMod1{key: 33}}
	weightMap = think.Export_makeWeightMap()
	tagMap = make(map[think.Module]bool)
	err = think.Export_calcWeight(weightMap, tms[11], tagMap)
	assert.Assert(t, err == nil, "calcWeight不检测依赖Module是否存在")
}
