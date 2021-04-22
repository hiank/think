package token

import (
	"context"
	"sync"
)

type Master struct {
	Token
	nextkey int            //NOTE: next key for cache token
	cache   map[int]*Slave //NOTE: derived token cache
	mux     sync.RWMutex
}

func NewMaster(ctx context.Context, uuid string) *Master {
	mtoken := &Master{
		cache: make(map[int]*Slave), //NOTE: only main Token need the cache
	}
	ctx, cancel := context.WithCancel(ctx)
	mtoken.Token = &baseToken{main: mtoken, uuid: uuid, Context: ctx, cancel: cancel}
	return mtoken
}

func (mtoken *Master) GetCached(key int) (tk Token, existed bool) {
	tk, existed = mtoken.cachedSlave(key)
	if existed && tk.Err() != nil {
		tk, existed = nil, false
		mtoken.deleteCached(key)
	}
	return
}

//ClearCache 清理缓存
func (mtoken *Master) ClearCache() {
	mtoken.mux.Lock()
	defer mtoken.mux.Unlock()
	for key, dtoken := range mtoken.cache {
		if dtoken.Err() != nil {
			delete(mtoken.cache, key)
		}
	}
}

func (mtoken *Master) cachedSlave(key int) (dtoken *Slave, existed bool) {
	mtoken.mux.RLock()
	defer mtoken.mux.RUnlock()
	dtoken, existed = mtoken.cache[key]
	return
}

func (mtoken *Master) deleteCached(key int) {
	mtoken.mux.Lock()
	defer mtoken.mux.Unlock()
	delete(mtoken.cache, key)
}

func (mtoken *Master) autoCache(dtoken *Slave) {
	mtoken.mux.Lock()
	defer mtoken.mux.Unlock()
	dtoken.key = mtoken.nextkey
	mtoken.cache[mtoken.nextkey] = dtoken
	mtoken.nextkey++
}

type Slave struct {
	Token
	key int
}

//Key key for cache
//key is 0 when the Token is main or not cached
func (dtoken *Slave) Key() int {
	return dtoken.key
}

type Token interface {
	context.Context
	Cancel()
	Derive(...Option) (Token, error)
	Uuid() string
}

type baseToken struct {
	context.Context
	cancel context.CancelFunc
	main   *Master //NOTE: main Token, would used by derived Token
	uuid   string
}

func (base *baseToken) Derive(opts ...Option) (Token, error) {
	dbase := &baseToken{uuid: base.uuid, main: base.main}
	dopts, dtoken := new(options), &Slave{Token: dbase}
	for _, opt := range opts {
		opt.apply(dopts)
	}

	if dopts.timeout > 0 {
		dbase.Context, dbase.cancel = context.WithTimeout(base.Context, dopts.timeout)
	} else {
		dbase.Context, dbase.cancel = context.WithCancel(base.Context)
	}
	if dopts.needCache {
		dbase.main.autoCache(dtoken)
	}
	return dtoken, nil
}

func (base *baseToken) Uuid() string {
	return base.uuid
}

func (base *baseToken) Cancel() {
	base.cancel()
}
