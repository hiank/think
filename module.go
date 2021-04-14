package think

import "context"

//Module for Launch
type Module interface {
	Depend() []Module               //NOTE: 依赖
	OnCreate(context.Context) error //NOTE: 初始化
	OnStart(context.Context) error  //NOTE: 启动
	OnStop()                        //NOTE: 停止阶段
	OnDestroy()                     //NOTE: 清理模块，程序结束前调用，用于落地数据或其它操作
}

//IgnoreDepend 忽略依赖方法
type IgnoreDepend int

//Depend 获取依赖
func (td IgnoreDepend) Depend() []Module {
	return []Module{}
}

//IgnoreOnCreate 忽略初始化过程
type IgnoreOnCreate int

//OnCreate 初始化
func (onc IgnoreOnCreate) OnCreate(context.Context) error {
	return nil
}

//IgnoreOnStart 忽略开始
type IgnoreOnStart int

//OnStart 开始
func (ons IgnoreOnStart) OnStart(context.Context) error {
	return nil
}

//IgnoreOnStop 忽略停止
type IgnoreOnStop int

//OnStop 停止
func (ons IgnoreOnStop) OnStop() {
}

//IgnoreOnDestroy 忽略释放
type IgnoreOnDestroy int

//OnDestroy 释放过程
func (ond IgnoreOnDestroy) OnDestroy() {
}
