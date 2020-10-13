package auth_test

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/hiank/think/auth"
)

func testTryGet(t *testing.T, key string, wait *sync.WaitGroup) {

	valCh, num, tkHub := make(chan *auth.Token, 10), 100, auth.NewTokenHub(context.Background())
	for i := 0; i < num; i++ {
		go func() {
			valCh <- tkHub.TryGet("111")
		}()
	}

	var lastVal *auth.Token
	var setOnce sync.Once
L:
	for {
		select {
		case val := <-valCh:
			setOnce.Do(func() {
				lastVal = val
			})
			assert.Equal(t, val.ToString(), lastVal.ToString(), "多goroutine 中，必须能正确获得val")
			if num--; num == 0 {
				break L
			}
		}
	}
	wait.Done()
}

func TestTokenHub(t *testing.T) {

	// auth.InitTokenHub(context.Background(), nil)
	tkHub := auth.NewTokenHub(context.Background())
	t.Run("TryGetPanic", func(t *testing.T) {
		defer func() {

			if r := recover(); r == nil {
				assert.Assert(t, false, "无效参数的话，必须抛出异常")
			}
		}()
		tkHub.TryGet("")
	})
	t.Run("TryGetSuccess", func(t *testing.T) {
		assert.Assert(t, tkHub.TryGet("111") != nil, "")
		assert.Assert(t, tkHub.TryGet("111") != tkHub.TryGet("111"), "获取是一个副本Token，多次获的Token是不一样的")
	})
	t.Run("TryGetMultiGoroutineSafe", func(t *testing.T) {

		wait, num := new(sync.WaitGroup), 100
		wait.Add(num)
		for i := 0; i < num; i++ {
			go testTryGet(t, strconv.Itoa(rand.Intn(100)), wait)
		}
		wait.Wait()
	})
}

func TestToken(t *testing.T) {

	tkHub := auth.NewTokenHub(context.Background())
	t.Run("必须能够派生token", func(t *testing.T) {
		tk := tkHub.TryGet("111")
		assert.Assert(t, tk.Derive() != nil, "Derive接口必须能获取一个新的token")
	})
	t.Run("派生token与原token不能为同一个", func(t *testing.T) {
		tk := tkHub.TryGet("111")
		assert.Assert(t, tk != tk.Derive(), "父token与派生的token不能是同一个")
	})
	t.Run("多次调用Derive接口得到派生token不能是同一个", func(t *testing.T) {
		tk := tkHub.TryGet("111")
		assert.Assert(t, tk.Derive() != tk.Derive(), "每次调用Derive得到的token不能是同一个")
	})
	t.Run("关闭将通知当前token及子token，不会影响父token及兄弟token", func(t *testing.T) {

		tkParent := tkHub.TryGet("111")
		tk, tkBrother := tkParent.Derive(), tkParent.Derive()
		tkChild := tk.Derive()
		tkGrandson := tkChild.Derive()

		tk.Invalidate()

		select {
		case <-tk.Done():
		default:
			t.Error("失效后，当前token需要收到Done通知")
		}

		select {
		case <-tkChild.Done():
		default:
			t.Error("失效后，子token需要收到Done通知")
		}

		select {
		case <-tkGrandson.Done():
		default:
			t.Error("失效后，孙token需要收到Done通知")
		}

		select {
		case <-tkBrother.Done():
			t.Error("失效后，兄弟token不能受影响")
		case <-tkParent.Done():
			t.Error("失效后，父token不能收影响")
		default:
		}
	})
	t.Run("only TokenHub can invalidate root token", func(t *testing.T) {

	})
	t.Run("所有派生线上的token有相同的ToString值", func(t *testing.T) {

		tkParent := tkHub.TryGet("111")
		tk, tkBrother := tkParent.Derive(), tkParent.Derive()
		tkChild := tk.Derive()
		tkGrandson := tkChild.Derive()

		assert.Equal(t, tkParent.ToString(), tk.ToString(), "父token与子token string值相等")
		assert.Equal(t, tkBrother.ToString(), tk.ToString(), "兄弟token string值相等")
		assert.Equal(t, tkParent.ToString(), tkGrandson.ToString(), "祖父子token间 string值相等")
	})
}

// func TestToken(t *testing.T) {

// 	auth.InitTokenHub(context.Background(), nil)
// 	t.Run("Invalidate", func(t *testing.T) {

// 		tk := auth.TryTokenHub().TryGet("111")
// 		tk.Invalidate()
// 		<-tk.Done()
// 	})
// 	t.Run("StrVal", func(t *testing.T) {

// 		tk := auth.TryTokenHub().TryGet("111")
// 		assert.Equal(t, tk.ToString(), "111", "")
// 	})
// 	t.Run("Timeout", func(t *testing.T) {
// 		tk := auth.TryTokenHub().TryGet("111")
// 		select {
// 		case <-tk.Done():
// 			assert.Assert(t, false, "必须过一段时间才会结束")
// 		default:
// 		}

// 		time.Sleep(time.Second)
// 		select {
// 		case <-tk.Done():
// 		default:
// 			assert.Assert(t, false, "超时后，token需要结束")
// 		}
// 	})
// 	auth.TryTokenHub().Close()
// }
