package set

import (
	"container/list"
	"encoding/json"
)

//loopUnmarshalJson 循环执行解析
func loopUnmarshalJson(list *list.List, val interface{}) (ok bool) {
	for iter := list.Front(); iter != nil; iter = iter.Next() {
		if json.Unmarshal(iter.Value.([]byte), val) == nil { //NOTE: 如果存在某一个解析出数据，则将返回true
			ok = true
		}
	}
	return
}
