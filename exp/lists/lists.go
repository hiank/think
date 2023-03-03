package lists

import "container/list"

func InsertBeforeFunc[T any](l *list.List, want T, firstEq func(cur, want T) bool) (out *list.Element) {
	var cut *list.Element
	for elm := l.Front(); elm != nil; elm = elm.Next() {
		if firstEq(elm.Value.(T), want) {
			cut = elm
			break
		}
	}
	if cut == nil {
		out = l.PushBack(want)
	} else {
		out = l.InsertBefore(want, cut)
	}
	return
}

func Foreach[T any](l *list.List, call func(val T) (done bool)) {
	for elm := l.Front(); elm != nil; elm = elm.Next() {
		if call(elm.Value.(T)) {
			break
		}
	}
}

func DeleteFunc[T any](l *list.List, cf func(val T) bool) {
	cur := l.Front()
	for cur != nil {
		tmp := cur
		cur = cur.Next()
		if cf(tmp.Value.(T)) {
			l.Remove(tmp)
		}
	}
}
