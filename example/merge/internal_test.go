package merge

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestConvertoDistcode(t *testing.T) {
	///
	// t.Error("unimplemented")

	sc := encodeSitecode(11, 1220)
	dc := convertoDistcode(sc)
	assert.Equal(t, dc.Layer(), uint8(11))
	assert.Equal(t, dc.baseIndex(), uint8(25))
	assert.DeepEqual(t, dc.Indexes(), []int{1220})

	assert.Equal(t, dc.Base(), encodeBaseDistcode(11, 25).Base())
}

func TestConvertoFiller(t *testing.T) {
	t.Error("unimplemented")
}

func TestRangeTwodimensional(t *testing.T) {
	// t.Error("unimplemented")
	var done bool
	rangeTwodimensional(nil, func(ia, ib, v int) { done = true })
	rangeTwodimensional([][]int{}, func(ia, ib, v int) { done = true })
	assert.Assert(t, !done, "no callback for empty slice")

	///
	rangeTwodimensional([][]int{
		{0, 1, 2},
		{1 << 16, (1 << 16) | 1, (1 << 16) | 2},
		{2 << 16, (2 << 16) | 1, (2 << 16) | 2},
	}, func(ia, ib, v int) {
		assert.Equal(t, (ia<<16)|ib, v)
		done = true
	})
	assert.Assert(t, done)
}

func TestFarmDataset(t *testing.T) {
	t.Error("unimplemented")
}

func TestCodes(t *testing.T) {
	t.Run("Sitecode", func(t *testing.T) {
		t.Run("exceed limit", func(t *testing.T) {
			// tag := 1
			var step int
			defer func(t *testing.T) {
				r := recover()
				assert.Equal(t, r, ErrInvalidValue)
				assert.Equal(t, step, 1)
			}(t)
			encodeSitecode(13, siteIndexLimit)
			step = 1
			encodeSitecode(13, siteIndexLimit+1)
			step = 2
		})
	})
}
