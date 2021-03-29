package pool

import (
	"context"
	"testing"

	testdatapb "github.com/hiank/think/net/testdata"
	"github.com/hiank/think/set/codes"
	"google.golang.org/protobuf/proto"
	"gotest.tools/v3/assert"
)

func TestLimitMuxThreadSafeRetain(t *testing.T) {

	maxLimit, wait, recvCnt := 1000, make(chan bool, 1000), 0
	limit := &LimitMux{max: maxLimit}

	for i := 0; i < maxLimit; i++ {
		go func() {
			limit.Retain()
			wait <- true
		}()
	}

	for range wait {
		recvCnt++
		if recvCnt == maxLimit {
			break
		}
	}
	assert.Equal(t, limit.cur, maxLimit)
}

func TestLimitMuxThreadSafeRelease(t *testing.T) {

	maxLimit, wait, recvCnt := 1000, make(chan bool, 1000), 0
	limit := &LimitMux{max: maxLimit, cur: maxLimit}

	for i := 0; i < maxLimit; i++ {
		go func() {
			limit.Release()
			wait <- true
		}()
	}

	for range wait {
		recvCnt++
		if recvCnt == maxLimit {
			break
		}
	}
	assert.Equal(t, limit.cur, 0)
}

func TestListMuxThreadSafePush(t *testing.T) {

	wait, cnt := make(chan int, 1000), 100000
	listMux := NewListMux()
	for i := 0; i < cnt; i++ {
		go func() {
			listMux.Push(&testdatapb.Test1{})
			wait <- 1
		}()
	}

	recvCnt := 0
	for range wait {
		recvCnt++
		if recvCnt == cnt {
			break
		}
	}
	assert.Equal(t, listMux.cache.Len(), cnt)
}

func TestListMuxThreadSafeShift(t *testing.T) {

	wait, cnt := make(chan int, 1000), 100000
	listMux := NewListMux()
	for i := 0; i < cnt; i++ {
		listMux.cache.PushBack(&testdatapb.Test1{})
	}
	for i := 0; i < cnt; i++ {
		go func() {
			listMux.Shift()
			wait <- 1
		}()
	}

	recvCnt := 0
	for range wait {
		recvCnt++
		if recvCnt == cnt {
			break
		}
	}
	assert.Equal(t, listMux.cache.Len(), 0)
}

func TestHubTrySetHandler(t *testing.T) {

	hub := NewHub(context.Background(), 100)

	t.Run("SetNilHandler", func(t *testing.T) {

		defer func() {
			code := recover()
			assert.Equal(t, code.(int), codes.PanicNilHandler, "if try set nil handler, funciton would panic PanicNilHandler")
		}()

		hub.TrySetHandler(nil)
	})
}

func TestHubThreadSafePushWithLaterHandler(t *testing.T) {
	hub, msgCnt, pushedCnt, pushed := NewHub(context.Background(), 100), 10000, 0, make(chan bool, 100)

	for i := 0; i < msgCnt; i++ {
		go func() {
			hub.Push(&testdatapb.Test1{})
			pushed <- true
		}()
	}

	for range pushed {
		pushedCnt++
		if pushedCnt == msgCnt {
			break
		}
	}
	assert.Equal(t, hub.list.cache.Len(), msgCnt)

	handleCh, handleCnt := make(chan bool, 100), 0
	hub.SetHandler(HandlerFunc(func(msg proto.Message) error {
		handleCh <- true
		return nil
	}))
	for range handleCh {
		handleCnt++
		if handleCnt == msgCnt {
			break
		}
	}
	assert.Equal(t, hub.list.cache.Len(), 0)
}

func TestHubThreadSafePushWithRandomHandler(t *testing.T) {

	hub, msgCnt, handleCnt, handleCh := NewHub(context.Background(), 100), 10000, 0, make(chan bool, 100)

	for i := 0; i < msgCnt; i++ {
		go hub.Push(&testdatapb.Test1{})
	}

	go func() {
		hub.SetHandler(HandlerFunc(func(msg proto.Message) error {
			handleCh <- true
			return nil
		}))
	}()

	for range handleCh {
		handleCnt++
		if handleCnt == msgCnt {
			break
		}
	}
	assert.Equal(t, hub.list.cache.Len(), 0)

}
