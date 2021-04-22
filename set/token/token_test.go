package token_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/hiank/think/set/token"
	"gotest.tools/v3/assert"
)

func TestDeriveDone(t *testing.T) {
	token := token.NewMaster(context.Background(), "BASETOKEN")
	derivedToken, err := token.Derive()
	assert.Assert(t, err == nil, err)

	token.Cancel()
	<-token.Done()
	assert.Equal(t, token.Err(), context.Canceled)

	<-derivedToken.Done()
	assert.Equal(t, derivedToken.Err(), context.Canceled)
	// assert.Assert(t, token.Err() == nil, token.Err())
}

func TestDeriveTimeout(t *testing.T) {
	tk := token.NewMaster(context.Background(), "BASETOKEN")
	derivedToken, err := tk.Derive(token.WithTimeout(time.Millisecond * 100))
	assert.Assert(t, err == nil, err)

	var ok bool
	select {
	case <-derivedToken.Done():
		ok = false
	case <-time.After(time.Millisecond * 90):
		ok = true
	}
	assert.Assert(t, ok, "超时时间未到，token不该Done")

	select {
	case <-time.After(time.Millisecond * 20):
		ok = false
	case <-derivedToken.Done():
		ok = true
	}
	assert.Assert(t, ok, "超时时间已过，token应该Done")
}

func TestDeriveDoneMainNotDone(t *testing.T) {
	tk := token.NewMaster(context.Background(), "BASETOKEN")
	derivedToken, err := tk.Derive() //token.WithTimeout(time.Millisecond * 100))
	assert.Assert(t, err == nil, err)

	derivedToken.Cancel()
	<-derivedToken.Done()
	assert.Equal(t, derivedToken.Err(), context.Canceled)

	assert.Assert(t, tk.Err() == nil, tk.Err())

	derivedToken, _ = tk.Derive(token.WithTimeout(time.Millisecond * 10))
	assert.Assert(t, derivedToken.Err() == nil, derivedToken.Err())

	derivedToken2, _ := derivedToken.Derive()
	assert.Assert(t, derivedToken2.Err() == nil, derivedToken2.Err())

	<-time.After(time.Millisecond * 20)
	assert.Equal(t, derivedToken.Err(), context.DeadlineExceeded)
	assert.Equal(t, derivedToken2.Err(), context.DeadlineExceeded)
	assert.Assert(t, tk.Err() == nil, tk.Err())
}

func TestDeriveDoneOtherNotDone(t *testing.T) {
	tk := token.NewMaster(context.Background(), "BASETOKEN")
	derivedToken, err := tk.Derive() //token.WithTimeout(time.Millisecond * 100))
	assert.Assert(t, err == nil, err)

	derivedToken2, _ := tk.Derive()
	assert.Assert(t, derivedToken2.Err() == nil, derivedToken2.Err())

	derivedToken.Cancel()
	<-derivedToken.Done()
	assert.Equal(t, derivedToken.Err(), context.Canceled)
	assert.Assert(t, derivedToken2.Err() == nil, derivedToken2.Err())
	assert.Assert(t, tk.Err() == nil, tk.Err())
}

func TestMasterCache(t *testing.T) {
	tk := token.NewMaster(context.Background(), "BASETOKEN")
	derivedToken, err := tk.Derive(token.WithCache()) //token.WithTimeout(time.Millisecond * 100))
	assert.Assert(t, err == nil, err)

	key := derivedToken.(*token.Slave).Key()
	dtoken, existed := tk.GetCached(key)
	assert.Assert(t, existed)
	assert.Equal(t, dtoken, derivedToken)

	_, existed = tk.GetCached(key + 1)
	assert.Assert(t, !existed)
}

func TestSyncSafe(t *testing.T) {
	tk := token.NewMaster(context.Background(), "BASETOKEN")
	// for
	// assert.Assert(t, false, "注意，当前没有考虑到，如果缓存的Token失效了，要清理这个缓存")
	dtoken, _ := tk.Derive(token.WithCache())
	cachedtk, _ := tk.GetCached(dtoken.(*token.Slave).Key())
	assert.Equal(t, dtoken, cachedtk)

	dtoken.Cancel()
	_, existed := tk.GetCached(dtoken.(*token.Slave).Key())
	assert.Assert(t, !existed)
}

func TestDeleteRepeat(t *testing.T) {
	hub := make(map[int]int)
	hub[1] = 11
	hub[2] = 12
	assert.Equal(t, len(hub), 2)

	delete(hub, 1)
	assert.Equal(t, len(hub), 1)

	delete(hub, 1)
	assert.Equal(t, len(hub), 1)
}

func BenchmarkDeleteRepeat(t *testing.B) {

	hub := make(map[int]int)
	hub[1] = 11
	for i := 0; i < 10000; i++ {
		delete(hub, 1)
	}
}

func BenchmarkDeleteWithJudge(t *testing.B) {

	hub := make(map[int]int)
	hub[1] = 11
	for i := 0; i < 10000; i++ {
		_, ok := hub[1]
		if ok {
			delete(hub, 1)
		}
	}
}

func TestContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var group sync.WaitGroup
	group.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			<-ctx.Done()
			group.Done()
		}()
	}
	<-time.After(time.Millisecond * 10)
	cancel()
	group.Wait()
	assert.Assert(t, true, "every goroutine can get ctx.Done")
}


func TestDeleteSafe(t *testing.T) {
	hub := make(map[int]int)
	hub[1] = 11
	hub[2] = 0
	hub[3] = 0
	hub[4] = 2
	hub[5] = 13
	hub[6] = 0
	for key, val := range hub {
		if val == 0 {
			delete(hub, key)
		}
	}
	assert.Equal(t, len(hub), 3)
	assert.Equal(t, hub[1], 11)
	assert.Equal(t, hub[4], 2)
	assert.Equal(t, hub[5], 13)
}
