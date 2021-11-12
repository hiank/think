package sys_test

import (
	"testing"

	"github.com/hiank/think/set/sys"
	"gotest.tools/v3/assert"
)

type testConf2 struct {
	Hostname string `json:"sys.hostname"`
	Hope     string `yaml:"Hope"`
	Hp       string `json:"sys.hp"`
	Slave    string `yaml:"redis.slave" json:"redis.slave"`
}

func TestHandleJson(t *testing.T) {
	tc, tc2 := &testConf{}, &testConf2{}
	sys.HandleJson("testdata/config.json", tc, tc2)
	sys.Unmarshal()
	assert.Equal(t, tc.Limit, 11)
	assert.Equal(t, tc2.Hostname, "hiank")
}

func TestHandleYaml(t *testing.T) {
	tc, tc2 := &testConf{}, &testConf2{}
	sys.HandleYaml("testdata/config.yaml", tc, tc2)
	sys.Unmarshal()
	assert.Equal(t, tc.Key, "ws")
	assert.Equal(t, tc2.Hope, "love")
}

func TestHandleFolder(t *testing.T) {
	tc, tc2 := &testConf{}, &testConf2{}
	sys.HandleFolder("testdata", tc, tc2)
	sys.Unmarshal()
	assert.Equal(t, tc.Key, "love-ws")
	assert.Equal(t, tc.Limit, 201)
	assert.Equal(t, tc2.Hope, "love")
	assert.Equal(t, tc2.Hostname, "hiank")
	assert.Equal(t, tc2.Hp, "hp")
	assert.Equal(t, tc2.Slave, "slave")

	t.Run("not folder param", func(t *testing.T) {
		// defer func(t *testing.T) {
		// 	r := recover()
		// 	assert.Assert(t, r != nil)
		// }(t)
		sys.HandleFolder("testdata/config.json", tc)
	})
}

// func TestHandleFolderDep(t *testing.T) {
// 	tc, tc2 := &testConf{}, &testConf2{}
// 	sys.HandleFolderDep("testdata", tc, tc2)
// 	sys.Unmarshal()
// 	assert.Equal(t, tc.Key, "love-ws")
// 	assert.Equal(t, tc.Limit, 201)
// 	assert.Equal(t, tc2.Hope, "love")
// 	assert.Equal(t, tc2.Hostname, "hiank")
// }
