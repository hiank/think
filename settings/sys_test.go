package settings

import (
	"testing"
	"github.com/hiank/conf"
	"gotest.tools/v3/assert"
)


func TestDefaultValue(t *testing.T) {

	assert.Equal(t, GetSys().WsPort, int16(8022))
	assert.Equal(t, GetSys().K8sPort, int16(8026))
	assert.Equal(t, GetSys().MessageGo, 1000)
	assert.Equal(t, GetSys().ClearInterval, 10)
}


func TestCustomValue(t *testing.T) {

	conf.LoadFromFile(GetSys(), "./settings_test.json")
	assert.Equal(t, GetSys().WsPort, int16(1024))
	assert.Equal(t, GetSys().K8sPort, int16(8026))
}