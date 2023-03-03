package play_test

import (
	"testing"

	"golang.org/x/exp/slices"
	"gotest.tools/v3/assert"
)

func TestMany(t *testing.T) {

	arr, nv := []int{1, 2, 4, 5}, 3
	idx := slices.IndexFunc(arr, func(v int) bool {
		return nv < v
	})
	arr = slices.Insert(arr, idx, nv)
	assert.DeepEqual(t, arr, []int{1, 2, 3, 4, 5})

	arr = slices.Insert(arr, len(arr), 6)
	assert.DeepEqual(t, arr, []int{1, 2, 3, 4, 5, 6})

	arr = arr[2:6]
	assert.DeepEqual(t, arr, []int{3, 4, 5, 6})
	// arr2 := []int{1, 2,3,4 }
}
