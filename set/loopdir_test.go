package set_test

import (
	"container/list"
	"encoding/json"
	"testing"

	"github.com/hiank/think/set"
	"gotest.tools/v3/assert"
)

type testJsonData struct {
	Value string `json:"value"`
}

func TestWalkText(t *testing.T) {
	cacheMap := make(map[string]*list.List)
	set.WalkText(cacheMap, "testdata/configs", "json")
	assert.Equal(t, len(cacheMap), 1)

	cache, ok := cacheMap["json"]
	assert.Assert(t, ok)

	assert.Equal(t, cache.Len(), 3, "测试目录下一共三个json结尾的文件")

	datas := make([]*testJsonData, 3)
	for iter, i := cache.Front(), 0; iter != nil; iter, i = iter.Next(), i+1 {
		val := &testJsonData{}
		err := json.Unmarshal(iter.Value.([]byte), val)
		assert.Assert(t, err == nil, err)
		datas[i] = val
	}

	cnt := 0
	wantDatas := []string{
		"value 1",
		"value 2",
		"value 3",
	}
L:
	for _, val := range wantDatas {
		for _, jsonVal := range datas {
			if jsonVal.Value == val {
				cnt++
				continue L
			}
		}
		assert.Assert(t, false, "测试数据中的值非期望值")
	}
	assert.Equal(t, cnt, 3)
}
