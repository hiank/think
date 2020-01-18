package setting

import (
	"testing"
	"gotest.tools/v3/assert"
)

type loaderTemp struct {

	Max int 
	Min int
}


func TestLoader(t *testing.T) {

	var temp loaderTemp
	err := LoadFromFile(&temp, "./setting_test.json")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, temp.Max, 10)
	assert.Equal(t, temp.Min, 5)
}