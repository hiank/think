package core

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"

	"gotest.tools/v3/assert"

	"github.com/golang/protobuf/ptypes/any"
)

type testMessage struct {
	tk  string
	val *any.Any
}

func (tm *testMessage) GetKey() string {
	return tm.tk
}

func (tm *testMessage) GetValue() *any.Any {
	return tm.val
}

func TestDoActive(t *testing.T) {

	mh := NewMessageHub(context.Background(), MessageHandlerTypeFunc(func(msg Message) error {
		return errors.New(msg.GetKey())
	}))

	rlt1, rlt2 := mh.Push(&testMessage{tk: "test", val: nil}), mh.Push(&testMessage{tk: "test2", val: nil})
	select {
	case <-rlt1:
		assert.Assert(t, false, "未激活的话，不能收到结果")
	case <-rlt2:
		assert.Assert(t, false, "未激活的话，不能收到结果")
	case <-time.After(time.Second):
	}

	mh.DoActive()
	cnt := 0

	for {
		select {
		case err1 := <-rlt1:
			assert.Equal(t, err1.Error(), "test", "收到处理结果")
			cnt++
		case err2 := <-rlt2:
			assert.Equal(t, err2.Error(), "test2", "收到处理结果")
			cnt++
		}
		if cnt == 2 {
			break
		}
	}
}

func TestSafeCall(t *testing.T) {

	numChan, stepChan := make(chan int), make(chan int)
	go func() {

		var num int
		for {
			if step, ok := <-stepChan; ok {
				num += step
			} else {
				break
			}
		}
		numChan <- num
	}()
	mh, wait := NewMessageHub(context.Background(), nil), new(sync.WaitGroup)
	for i := 0; i < 10000; i++ {
		switch rand.Intn(4) {
		case 0:
			fallthrough
		case 1:
			wait.Add(1)
			go func() {
				mh.safePushBack(new(messageReq))
				stepChan <- 1
				wait.Done()
			}()
		case 2:
			wait.Add(1)
			go func() {
				if mh.safeShift() != nil {
					stepChan <- -1
				}
				wait.Done()
			}()
		case 3:
			// t.Log(mh.safeLen())
		}
	}
	wait.Wait()
	close(stepChan)
	// time.Sleep(time.Millisecond * 10)
	assert.Equal(t, <-numChan, mh.safeLen(), "测试通过，没有panic，并且数量正确")
}
