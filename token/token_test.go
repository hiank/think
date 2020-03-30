package token

import (
	"container/list"
	"context"
	"sync"
	"testing"
	"time"
	"gotest.tools/v3/assert"
)

func TestMapDelete(t *testing.T) {

	m := make(map[string]int)

	m["1"] = 10

	delete(m, "1")
	delete(m, "2")		//NOTE: 测试删除不存在的元素
	t.Log(m)
}

func TestOnceDo(t *testing.T) {

	var once sync.Once
	num := 0
	onceFunc := func ()  {
		num++
	}
	once.Do(onceFunc)
	once.Do(onceFunc)
	assert.Equal(t, num, 1)

	num2 := 0
	once.Do(func ()  {
	
		num++
		num2++
	})
	assert.Equal(t, num2, 0)
	assert.Equal(t, num, 1)
}


func TestAsyncOnceDo(t *testing.T) {

	num := 0
	onceFunc := func ()  {
		
		time.Sleep(1000)
		num++
	}

	ch := make(chan int, 2)
	var once sync.Once

	go func ()  {
		
		once.Do(onceFunc)
		t.Log("1")
		assert.Equal(t, num, 1)
		ch <- 1
	}()

	go func () {
		
		once.Do(onceFunc)
		t.Log("2")
		assert.Equal(t, num, 1)
		ch <- 2
	}()

	<- ch
	<- ch
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
	doneFunc := func ()  {
		
		<-ctx.Done()
		num++
		wait.Done()
	}
	go doneFunc()
	go doneFunc()

	go func ()  {
		<-time.After(time.Second)
		cancel()
	}()
	// <-time.After(time.Second)
	// cancel()
	wait.Wait()
	assert.Equal(t, num, 2)
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

func TestListRemove(t *testing.T) {

	queue := list.New()
	element := queue.PushBack(1)
	queue.Init()
	queue.Remove(element)		//NOTE: 测试表明，经过Init 之后，某个已知的element 状态不会重置，导致再删的时候，会错乱

	queue.PushBack(2)
	assert.Equal(t, queue.Len(), 0)
}


func TestTokenCancel(t *testing.T) {

	tk, _ := newToken(context.WithValue(context.Background(), IdentityKey, 1001))
	derivedToken := tk.Derive()
	// assert.Equal(t, derivedToken.Value(contextKey(1)).(int), 1002)
	assert.Equal(t, derivedToken.Value(IdentityKey).(int), 1001)

	// assert.Equal(t, tk.queue.Len(), 1)
	// tk.Cancel()

	// assert.Equal(t, tk.queue.Len(), 1)		//NOTE: 监听goroutine 中将删除元素，所以此时并没有立即删除
	// time.Sleep(time.Microsecond)

	// assert.Equal(t, tk.queue.Len(), 0)		//NOTE: 等待一定时间后，元素被删除
}

// func TestBuilder(t *testing.T) {

// 	ctx, cancel := context.WithCancel(context.Background())
// 	token.InitBuilder(ctx)
// 	defer token.ReleaseBuilder()
// 	tokenObj, _, err := token.Get("2001")
// 	if err != nil {
// 		t.Log(err)
// 	}

// 	t.Log(tokenObj.ToString())
// 	go func() {

// 		tokenObj.Cancel()
// 		// t.Log()
// 		<-time.After(time.Second)
// 		cancel()
// 	}()

// 	<-ctx.Done()
// 	t.Log("passed test")
// }