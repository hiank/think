package run

import (
	"context"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestTaskWorker(t *testing.T) {
	// wc := make(chan Task)
	tl := newTaskWorker(context.Background(), 0) //&taskList{l: list.New(), wc: wc}
	// tl.wc = make(chan Task)

	assert.Equal(t, tl.wt, nil)
	// assert.Assert(t, tl.wc != nil)
	// tl.
	sc, t1 := tl.work()
	assert.Assert(t, sc == nil)
	assert.Equal(t, t1, nil)

	t2 := NewLiteTask(func(int) error {
		<-time.After(time.Millisecond * 10)
		return nil
	}, 1)
	tl.push(t2)
	assert.Equal(t, tl.wt, nil, "first task sended")

	for i := 0; i < cap(tl.wc); i++ {
		tl.push(NewLiteTask(func(t int) error {
			<-time.After(time.Millisecond * 10)
			return nil
		}, 2))
	}
	assert.Equal(t, cap(tl.wc), len(tl.wc), "cache full")

	t3 := NewLiteTask(func(int) error {
		<-time.After(time.Millisecond * 100)
		return nil
	}, 3)
	tl.push(t3)
	assert.Equal(t, tl.l.Len(), 0, "注意，这个测试时间性，task返回过快的化可能导致测试失败. 但是很罕见")
	assert.Equal(t, tl.wt, t3, "wait for send to work chan")

	t4 := NewLiteTask(func(int) error { return nil }, 4)
	tl.push(t4)
	assert.Equal(t, tl.l.Len(), 1)

	sc, t1 = tl.work()
	assert.Assert(t, sc != nil)
	assert.Equal(t, t1, t3)
	assert.Equal(t, tl.l.Len(), 1)

	tl.wt = nil
	sc, t1 = tl.work()
	assert.Assert(t, sc != nil)
	assert.Equal(t, t1, t4)
	assert.Equal(t, tl.l.Len(), 0)
}
