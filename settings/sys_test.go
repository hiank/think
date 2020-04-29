package settings

import (
	"testing"

	"github.com/hiank/conf"
	"gotest.tools/v3/assert"
)

func TestDefaultValue(t *testing.T) {

	assert.Equal(t, GetSys().WsPort, uint16(8022))
	assert.Equal(t, GetSys().K8sPort, uint16(8026))
	assert.Equal(t, GetSys().MessageGo, 1000)
	assert.Equal(t, GetSys().ClearInterval, 10)
}

func TestCustomValue(t *testing.T) {

	conf.LoadFromFile(GetSys(), "./sys_test.json")
	assert.Equal(t, GetSys().WsPort, uint16(1024))
	assert.Equal(t, GetSys().K8sPort, uint16(8026))
}
