package lists_test

import (
	"container/list"
	"testing"

	"github.com/hiank/think/exp/lists"
	"gotest.tools/v3/assert"
)

func TestInsertBeforeFunc(t *testing.T) {
	l := list.New()
	firstEq := func(cur, want int) bool {
		return cur > want
	}
	// elm := l.InsertBefore(1, l.Front())
	// assert.Equal(t, elm.Value.(int), 1)
	elm := lists.InsertBeforeFunc(l, 11, firstEq)
	assert.Equal(t, elm.Value.(int), 11)

	elm = lists.InsertBeforeFunc(l, 22, firstEq)
	assert.Equal(t, elm.Value.(int), 22)

	assert.Equal(t, l.Front().Value.(int), 11)
	assert.Equal(t, l.Back().Value.(int), 22)

	lists.InsertBeforeFunc(l, 13, firstEq)
	assert.Equal(t, l.Front().Value.(int), 11)
	assert.Equal(t, l.Front().Next().Value.(int), 13)
	assert.Equal(t, l.Back().Value.(int), 22)

	// lists.InsertBeforeFunc(l, 13, firstEq)
}

type tmpStruct struct {
	v int
}

func TestForeach(t *testing.T) {
	l := list.New()
	l.PushBack(&tmpStruct{v: 1})
	l.PushBack(&tmpStruct{v: 3})
	var v *tmpStruct
	l.PushBack(v)
	l.PushBack(&tmpStruct{v: 2})

	var bit = 0
	lists.Foreach(l, func(v *tmpStruct) (done bool) {
		if done = v == nil; !done {
			bit |= 1 << v.v
		}
		// bit |= 1 << v
		return
	})

	assert.Equal(t, bit, (1<<1)|(1<<3))
}

func TestDeleteFunc(t *testing.T) {
	l := list.New()
	l.PushBack(1)
	l.PushBack(3)
	l.PushBack(2)
	l.PushBack(7)
	l.PushBack(5)

	lists.DeleteFunc(l, func(v int) bool {
		return v > 2
	})
	assert.Equal(t, l.Len(), 2)
	assert.Equal(t, l.Front().Value.(int), 1)
	assert.Equal(t, l.Back().Value.(int), 2)
}
