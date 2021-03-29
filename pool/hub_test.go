package pool_test

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/hiank/think/net/pb"
	testdatapb "github.com/hiank/think/net/testdata"
	"github.com/hiank/think/pool"
	"gotest.tools/v3/assert"
)

// func TestLimitMux(t *testing.T) {

// 	lm := &pool.LimitMux{Max: 1}
// 	assert.Assert(t, !lm.Locked(), "未到限制数时，不锁定")
// 	assert.Assert(t, lm.Retain(), "未到限制数时，可请求")
// 	assert.Assert(t, lm.Locked(), "到限制数，锁定")
// 	assert.Assert(t, !lm.Retain(), "到限制数，不能请求")
// }

func TestListMux(t *testing.T) {

	listMux := pool.NewListMux()
	assert.Equal(t, listMux.Shift(), nil, "没有数据时，Shift返回nil")
	listMux.Push(&pb.Message{Key: "test1"})
	listMux.Push(&pb.Message{Key: "test2"})
	assert.Equal(t, listMux.Shift().(*pb.Message).GetKey(), "test1", "有数据时，Shift返回最前面数据")
	assert.Equal(t, listMux.Shift().(*pb.Message).GetKey(), "test2", "")
	assert.Equal(t, listMux.Shift(), nil, "")
}

func TestHubPoolRemoveAsyncSafe(t *testing.T) {

	hp, max := pool.NewHubPool(context.Background()), 10000
	for i := 0; i < max; i++ {
		hp.AutoHub(strconv.Itoa(i))
	}

	for i := 0; i < max; i++ {
		hub := hp.GetHub(strconv.Itoa(i))
		assert.Assert(t, hub != nil)
	}

	ch := make(chan bool, 1000)
	for i := 0; i < max; i++ {
		go func(i int) {
			hp.Remove(strconv.Itoa(i))
			ch <- true
		}(i)
	}

	for i := 0; i < max; i++ {
		go func() {
			hp.Remove(strconv.Itoa(rand.Intn(max))) //NOTE: 测试重复删除安全性
			ch <- true
		}()
	}

	recvCnt := 0
	for range ch {
		recvCnt++
		if recvCnt == max*2 {
			break
		}
	}

	for i := 0; i < max; i++ {
		hub := hp.GetHub(strconv.Itoa(i))
		assert.Assert(t, hub == nil, "删除所有内容后，hub将不再包含")
	}
}

func TestHubPoolRemoveAllAsyncSafe(t *testing.T) {

	hp, max := pool.NewHubPool(context.Background()), 10000
	for i := 0; i < max; i++ {
		hp.AutoHub(strconv.Itoa(i))
	}

	for i := 0; i < max; i++ {
		hub := hp.GetHub(strconv.Itoa(i))
		assert.Assert(t, hub != nil)
	}

	ch := make(chan bool, 1000)
	for i := 0; i < max; i++ {
		go func() {
			hp.RemoveAll()
			ch <- true
		}()
	}

	recvCnt := 0
	for range ch {
		recvCnt++
		if recvCnt == max {
			break
		}
	}

	for i := 0; i < max; i++ {
		hub := hp.GetHub(strconv.Itoa(i))
		assert.Assert(t, hub == nil, "删除所有内容后，hub将不再包含")
	}
}

//BenchmarkHubPoolAsync
func BenchmarkHubPoolAsync(t *testing.B) {

	hp, max := pool.NewHubPool(context.Background()), 1000
	cache, ch := make(map[int]int), make(chan int, 1000)

	for i := 0; i < max; i++ {
		_, ok := cache[i]
		assert.Assert(t, !ok)
	}
	for x := 0; x < max; x++ {
		for i := 0; i < max; i++ {
			go func(x int) {
				if _, isNew := hp.AutoHub(strconv.Itoa(x)); isNew {
					ch <- x
				} else {
					ch <- -1
				}
			}(x)
		}
	}

	recvCnt := 0
	for val := range ch {
		recvCnt++
		if val >= 0 {
			lastVal, ok := cache[val]
			assert.Assert(t, !ok, fmt.Sprintf("不能重复赋值key:%d, val:%d", val, lastVal))
			cache[val] = val
		}
		if recvCnt == max*max {
			break
		}
	}
}

//BenchmarkListMuxAsync 延时ListMux是线程安全的
func BenchmarkListMuxAsync(t *testing.B) {

	wait, cnt := make(chan int, 1000), 100000
	listMux := pool.NewListMux()
	go func() {

		for i := 0; i < cnt; i++ {
			go func() {
				listMux.Push(&testdatapb.Test1{})
				wait <- 1
			}()
		}
	}()

	go func() {
		for i := 0; i < cnt; i++ {
			go func() {
				listMux.Shift()
				wait <- 1
			}()
		}
	}()

	recvCnt := 0
	for range wait {
		recvCnt++
		if recvCnt == cnt*2 {
			break
		}
	}
}
