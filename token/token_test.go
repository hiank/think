package token

import (
	"container/list"
	"context"
	"sync"
	"testing"
	"time"

	"github.com/hiank/think/settings"
	"gotest.tools/v3/assert"
)

func TestMapDelete(t *testing.T) {

	m := make(map[string]int)

	m["1"] = 10

	delete(m, "1")
	delete(m, "2") //NOTE: 测试删除不存在的元素
	t.Log(m)
}

func TestOnceDo(t *testing.T) {

	var once sync.Once
	num := 0
	onceFunc := func() {
		num++
	}
	once.Do(onceFunc)
	once.Do(onceFunc)
	assert.Equal(t, num, 1)

	num2 := 0
	once.Do(func() {

		num++
		num2++
	})
	assert.Equal(t, num2, 0)
	assert.Equal(t, num, 1)
}

func TestAsyncOnceDo(t *testing.T) {

	num := 0
	onceFunc := func() {

		time.Sleep(1000)
		num++
	}

	ch := make(chan int, 2)
	var once sync.Once

	go func() {

		once.Do(onceFunc)
		t.Log("1")
		assert.Equal(t, num, 1)
		ch <- 1
	}()

	go func() {

		once.Do(onceFunc)
		t.Log("2")
		assert.Equal(t, num, 1)
		ch <- 2
	}()

	<-ch
	<-ch
	t.Log("3")
}

func TestContextValue(t *testing.T) {

	ctx := context.WithValue(context.Background(), "level", 1)
	ctx = context.WithValue(ctx, "level", 2)
	lv := ctx.Value("level").(int)
	t.Log(lv)

	val := ctx.Value(IdentityKey)
	t.Log(val)
}

//TestContextDone 用于验证，当context Cancel 调用后，所有此context 的Done() 都会响应
func TestContextDone(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	num, wait := 0, new(sync.WaitGroup)
	wait.Add(2)
	doneFunc := func() {

		<-ctx.Done()
		num++
		wait.Done()
	}
	go doneFunc()
	go doneFunc()

	go func() {
		<-time.After(time.Second)
		cancel()
	}()

	wait.Wait()
	assert.Equal(t, num, 2)
}

func TestConextCancel(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	// ctx1, cancel1 := context.WithCancel(ctx)
	wait := new(sync.WaitGroup)
	wait.Add(1)
	go func(ctx context.Context) {

		ctx1, cancel1 := context.WithCancel(ctx)
		select {
		case <-ctx1.Done():
			t.Log("2", ctx1.Done())
			cancel1()
			cancel1() //NOTE: 经过测试，cancel 可以多次调用。父Context关闭，子Context 可以接收到Done
			// case <-ctx.Done():
			// 	t.Log("1", ctx.Err())
			// 	cancel1()
		}
		wait.Done()
	}(ctx)

	<-time.After(time.Second)

	cancel()
	wait.Wait()
}

func TestListRemove(t *testing.T) {

	queue := list.New()
	element := queue.PushBack(1)
	queue.Init()
	queue.Remove(element) //NOTE: 测试表明，经过Init 之后，某个已知的element 状态不会重置，导致再删的时候，会错乱

	queue.PushBack(2)
	assert.Equal(t, queue.Len(), 0)
}

func TestTokenCancel(t *testing.T) {

	tk, _ := newToken(context.WithValue(context.Background(), IdentityKey, "1001"))
	derivedToken := tk.Derive()
	assert.Equal(t, derivedToken.Value(IdentityKey).(string), "1001")

	tok := GetBuilder().Get("test")
	tok.Cancel()
	tok1, ok := GetBuilder().Find("test")
	assert.Equal(t, ok, false)
	assert.Equal(t, tok1, nilToken)
	// assert.Equal(t, )
}

func TestTokenDerive(t *testing.T) {

	tok := GetBuilder().Get("test")
	tok1 := tok.Derive()
	tok1.Cancel()
	tok2, ok := GetBuilder().Find("test") //NOTE: 此处表明派生的token关闭后并不影响父token的状态
	assert.Equal(t, tok2, tok)
	assert.Equal(t, ok, true)

	tok1 = tok.Derive()
	go func() {
		tok.Cancel()
	}()
	select {
	case <-tok1.Done():
		assert.Assert(t, true) //NOTE: 此处表明，父token关闭后，派生的token也会收到Done消息
	}
}

//基于定时器的测试可能发生概率性的报错，这个注意一下
func TestTokenTimeout(t *testing.T) {

	settings.GetSys().TimeOut = 100
	GetBuilder().Get("test")
	// tok := GetBuilder().
	<-time.After(time.Millisecond * 80)
	tok1, ok := GetBuilder().Find("test")
	assert.Assert(t, ok)
	assert.Assert(t, tok1 != nilToken)

	<-time.After(time.Millisecond * 40)

	tok1, ok = GetBuilder().Find("test")
	assert.Assert(t, !ok)
	assert.Assert(t, tok1 == nilToken)
}

// 这不是一个稳定的测试，因为基于定时器，容易出现误差，导致偶发的错误
// 当需要测试定时器逻辑时，可打开这个方法单独测试，其余时间建议关闭，避免偶发性的报错
// 当前错误的发生可能是与时间量级相关的，当前测试时间是微米级的，如果扩大为s级别，偶发性的错误应该不再出现[未验证]
// func TestTokenResetTimer(t *testing.T) {

// 	settings.GetSys().TimeOut = 100
// 	GetBuilder().Get("test1")
	
// 	<-time.After(time.Millisecond * 90)

// 	tok, ok := GetBuilder().Find("test1")
// 	assert.Assert(t, ok)

// 	tok.ResetTimer()
// 	<-time.After(time.Millisecond * 110)
// 	tok1, ok := GetBuilder().Find("test1")
// 	assert.Assert(t, !ok)
// 	assert.Equal(t, tok1, nilToken)

// 	tok1 = GetBuilder().Get("test1")
// 	<-time.After(time.Millisecond * 50)
// 	for i:=0; i<100; i++ {
// 		go tok1.ResetTimer()
// 	}
// 	<-time.After(time.Millisecond * 90)
// 	tok2, ok := GetBuilder().Find("test1")
// 	assert.Assert(t, ok)
// 	assert.Equal(t, tok1, tok2)

// 	<-time.After(time.Millisecond * 20)
// 	tok2, ok = GetBuilder().Find("test1")
// 	assert.Assert(t, !ok)
// 	assert.Equal(t, tok2, nilToken)
// }
