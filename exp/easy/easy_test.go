package easy_test

import (
	"strconv"
	"testing"

	"github.com/hiank/think/exp/easy"
	"gotest.tools/v3/assert"
)

type tmpIntoString struct {
}

func (tis *tmpIntoString) Convert(v int) string {
	return strconv.FormatInt(int64(v), 10)
}

func TestTowDimensional(t *testing.T) {
	v1 := [][]int{
		{11, 12, 13},
		{21, 22, 23},
	}
	// var v2 [][]string
	v2 := easy.TwoDimensinal[int, string](v1, &tmpIntoString{})
	assert.DeepEqual(t, v2, [][]string{
		{"11", "12", "13"},
		{"21", "22", "23"},
	})
}

func TestResetBit(t *testing.T) {
	i, _ := strconv.ParseInt("1001101001000", 2, 64)
	i = easy.ResetBit(i, 0, 4, 8)
	str := strconv.FormatInt(i, 2)
	assert.Equal(t, str, "1000000001000")

	i8 := easy.ResetBit(int8(-1), 0, 0, 8)
	assert.Equal(t, i8, int8(0))

	i8 = easy.ResetBit(int8(-1), 0, 7, 1)
	assert.Equal(t, i8, int8(127))

	i8 = easy.ResetBit(int8(-1), 2, 2, 4)
	//convert "11001011" -> "11001010" -> "10110101"
	i82, _ := strconv.ParseInt("-110101", 2, 8)
	assert.Equal(t, i8, int8(i82))
}

func TestBitValue(t *testing.T) {
	i, _ := strconv.ParseInt("1001101001000", 2, 64)
	v := easy.BitValue(i, 3, 1)
	assert.Equal(t, v, int64(1))

	assert.Equal(t, easy.BitValue(i, 5, 4), int64(10))
}

// func TestMake(t *testing.T) {
// 	sli, err := easy.Make[[]int](1,2,3)
// 	assert.Equal(t, err, nil)
// 	assert.Equal(t, len(sli), 1)
// 	assert.Equal(t, cap(sli), 2)

// 	m, err := easy.Make[map[]]()
// }

func TestInstantiate(t *testing.T) {
	v, err := easy.Instantiate[int]()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 0)

	v1, err := easy.Instantiate[*int]()
	assert.Equal(t, err, nil)
	assert.Equal(t, *v1, 0)

	_, err = easy.Instantiate[[]int]()
	assert.Equal(t, err, easy.ErrUnsupportType)
}
