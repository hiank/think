package doc_test

import (
	"reflect"
	"testing"

	"github.com/hiank/think/doc"
)

const loopCnt = 100000

func BenchmarkMakeT(b *testing.B) {
	for i := 0; i < loopCnt; i++ {
		doc.MakeT[[]tmpExcel](1024, 1024)
		doc.MakeT[map[string]tmpExcel]()
		doc.MakeT[tmpExcel]()
		doc.MakeT[*tmpExcel]()
	}
}

func BenchmarkReflectNew(b *testing.B) {
	for i := 0; i < loopCnt; i++ {
		var s []tmpExcel
		rt := reflect.TypeOf(s)
		reflect.MakeSlice(rt, 1024, 1024)

		var m map[string]tmpExcel
		rt = reflect.TypeOf(m)
		reflect.MakeMap(rt)

		var v tmpExcel
		reflect.New(reflect.TypeOf(v))

		var v2 *tmpExcel
		reflect.New(reflect.TypeOf(v2).Elem())
	}
}
