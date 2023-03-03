package easy

import "github.com/hiank/think/run"

const (
	ErrExceedBitLimit = run.Err("exceed the bit limit")
)

type Integer interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~int8 | ~int16 | ~int32 | ~int64 | ~int
}

func ResetBit[T Integer](v, bitV T, low, cnt uint) T {
	var bit T = ((T(1) << cnt) - 1) << low
	tmp := bitV << low
	if bitV = tmp & bit; bitV != tmp {
		panic(ErrExceedBitLimit)
	}
	var max T = T(0) - 1 //111111...
	return ((max ^ bit) & v) | bitV
}

func BitValue[T Integer](v T, low, cnt uint) (bitV T) {
	bitV = v >> low
	var cover T = T(1<<cnt) - 1
	bitV &= cover
	return
}

// func SliceInsertFunc[S ~[]E, E any](s S, f func(E) bool, v E) S {
// 	idx := slices.IndexFunc(s, f)
// 	if idx == -1 {
// 		idx = len(s)
// 	}
// 	return slices.Insert(s, idx, v)
// }
