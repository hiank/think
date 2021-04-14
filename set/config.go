package set

import (
	"container/list"
	"reflect"

	"k8s.io/klog/v2"
)

var (
	JSON = "json"
	YAML = "yaml"
)

func UnmarshalJSON(textlist *list.List, valist *list.List) {
	for valist.Len() > 0 {
		val := valist.Remove(valist.Front())
		if !loopUnmarshalJson(textlist, val) {
			klog.Warningf("cannot unmarshal %v from config files\n", reflect.TypeOf(val))
		}
	}
}
