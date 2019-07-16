package token_test

import (
	"github.com/hiank/think/token"
	"sync"
	"time"
	"context"
	"testing"
)

func TestMapDelete(t *testing.T) {

	m := make(map[string]int)

	m["1"] = 10

	delete(m, "1")
	delete(m, "2")		//NOTE: 测试删除不存在的元素
	t.Log(m)
}


func TestContextValue(t *testing.T) {

	ctx := context.WithValue(context.Background(), "level", 1)
	ctx = context.WithValue(ctx, "level", 2)
	lv := ctx.Value("level").(int)
	t.Log(lv)

	val := ctx.Value(token.ContextKey("token"))
	t.Log(val)
}

func TestConextCancel(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	// ctx1, cancel1 := context.WithCancel(ctx)
	wait := new(sync.WaitGroup)
	wait.Add(1)
	go func (ctx context.Context) {

		ctx1, cancel1 := context.WithCancel(ctx)
		select {
		case <-ctx1.Done():
			t.Log("2", ctx1.Done())
			cancel1()
			cancel1()				//NOTE: 经过测试，cancel 可以多次调用。父Context关闭，子Context 可以接收到Done
		// case <-ctx.Done():
		// 	t.Log("1", ctx.Err())
		// 	cancel1()
		}
		wait.Done()
	} (ctx)

	<-time.After(time.Second)

	cancel()
	wait.Wait()
}

func TestBuilder(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	token.InitBuilder(ctx)
	tokenObj, err := token.Get("2001")
	if err != nil {
		t.Log(err)
	}

	t.Log(tokenObj.ToString())
	go func() {

		tokenObj.Cancel()
		// t.Log()
		<-time.After(time.Second)
		cancel()
	}()

	<-ctx.Done()
	t.Log("passed test")
}