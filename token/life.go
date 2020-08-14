package token

import (
	"sync"
	"context"
)

// Life 用于
type Life struct {
	context.Context
	Kill context.CancelFunc
}

// Derive 派生新的Life
func (life *Life) Derive() *Life {

	derived := new(Life)
	derived.Context, derived.Kill = context.WithCancel(life.Context)
	return derived
}


var _singleLife *Life
var _singleLifeOnce sync.Once

//BackgroundLife 获得Life
func BackgroundLife() *Life {

	_singleLifeOnce.Do(func() {

		_singleLife = new(Life)
		_singleLife.Context, _singleLife.Kill = context.WithCancel(context.Background())
		go func() {
			<-_singleLife.Done()
			_singleLife = nil
			_singleLifeOnce = sync.Once{}
		}()
	})
	return _singleLife
}