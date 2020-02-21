package token

import (
	"testing"
	"time"
	"gotest.tools/v3/assert"
	// is "gotest.tools/assert/cmp"
)



func TestGet(t *testing.T) {

	_, err := GetBuilder().Get("TestGet")
	if err != nil {
		assert.Assert(t, false)
		return
	}
	// assert.Assert(t, false)
}


func TestFind(t *testing.T) {

	_, ok := GetBuilder().Find("TestFind")
	assert.Assert(t, !ok)

	GetBuilder().Build("TestFind")
	_, ok = GetBuilder().Find("TestFind")
	assert.Assert(t, ok)
}

func TestBuild(t *testing.T) {

	_, err := GetBuilder().Build("TestBuild")
	assert.Assert(t, err == nil)
	_, err = GetBuilder().Build("TestBuild")
	assert.Assert(t, err != nil)
}

func TestDelete(t *testing.T) {

	GetBuilder().Build("TestDelete")
	_, ok := GetBuilder().Find("TestDelete")
	assert.Assert(t, ok)
	GetBuilder().Delete("TestDelete")
	_, ok = GetBuilder().Find("TestDelete")
	assert.Assert(t, !ok)
}


//NOTE: 测试单例调用
func TestCancel(t *testing.T) {

	var nilVal *Builder

	assert.Assert(t, GetBuilder() != nilVal)
	GetBuilder().Cancel()

	time.Sleep(1000)			//NOTE: 需要等待监听goroutine 处理ctx.Done()
	assert.Equal(t, GetBuilder(), nilVal)
}