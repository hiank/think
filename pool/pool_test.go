package pool


import (
	"container/list"
	"testing"

	"gotest.tools/assert"
	// "github.com/stretchr/testify/assert"
)

//测试 list.List 包Remove 对
func TestListRemove(t *testing.T) {

	queue := list.New()
	// queue.PushBack(1)
	for i:=0; i<10; i++ {
		queue.PushBack(i)
	}

	num := 0
	element:=queue.Front()
	for element != nil {

		tmp := element
		element = element.Next()
		// element.Value
		queue.Remove(tmp)
		num++
	}
	assert.Equal(t, num, 10)
}

//测试list Remove，当element 被remove后，迭代将失效
func TestListRemove2(t *testing.T) {

	queue := list.New()
	// queue.PushBack(1)
	for i:=0; i<10; i++ {
		queue.PushBack(i)
	}

	element := queue.Front()
	queue.Remove(element)
	var nilObj *list.Element
	assert.Equal(t, element.Next(), nilObj)
}


// func TestContextDone(t *testing.T) {

// 	ctx := context.Background()
// 	// ctx.Done()
// 	// context.WithTimeout(1)
// 	go func() {
		
// 		select {
// 		case <- ctx.Done():
// 			t.Log("recive ctx done")
// 		}
// 	}()

// 	t.Log("before close ctx")
// 	// close(ctx.Done())
// 	// close(ctx)
// 	// ctx.Done() <- true
// }


// func TestPoolClose(t *testing.T) {


// }