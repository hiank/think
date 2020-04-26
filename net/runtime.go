package net

import (
	"context"
	"sync"
)

//Runtime 运行时
type Runtime struct {

	context.Context
	Close 	context.CancelFunc
}


var _singleRuntime *Runtime
var _singleRuntimeOnce sync.Once

//GetRuntime 获得Runtime
func GetRuntime() *Runtime {

	_singleRuntimeOnce.Do(func ()  {
		
		_singleRuntime = new(Runtime)
		_singleRuntime.Context, _singleRuntime.Close = context.WithCancel(context.Background())
		go func ()  {
			<-_singleRuntime.Done()
			_singleRuntime = nil
			_singleRuntimeOnce = sync.Once{}
		}()
	})
	return _singleRuntime
}



