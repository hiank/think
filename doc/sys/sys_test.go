package sys_test

import (
	"testing"

	"github.com/hiank/think/doc/sys"
	"gotest.tools/v3/assert"
)

type tmpJson struct {
	Id       int    `json:"tmp.id"`
	Name     string `json:"tmp.name"`
	Limit    uint   `json:"sys.Limit"`
	Hostname string `json:"sys.hostname"`
}

type tmpYaml struct {
	Key   string `yaml:"ws.k"`
	Value int    `yaml:"ws.v"`
	A     string `yaml:"a"`
	Hope  string `yaml:"Hope"`
	Hp    string `yaml:"hp"`
	Slave string `yaml:"redis.slave"`
}

const tmpJsonValue = `{
	"tmp.id": 9527,
	"tmp.name": "华安"
}`
const tmpYamlValue = `
ws.k: "bbq"
ws.v: 25
`

func TestFat(t *testing.T) {
	fat := sys.NewFat()
	fat.Load("testdata/dep")
	var jv tmpJson
	var yv tmpYaml
	err := fat.UnmarshalTo(&jv, &yv)
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, jv, tmpJson{Limit: 1})
	assert.DeepEqual(t, yv, tmpYaml{A: "love-ws"})

	fat.Load("testdata")
	fat.UnmarshalTo(&jv, &yv)
	//testdata/dep 不会再加载
	//testdata/dep2 最后加载
	assert.DeepEqual(t, jv, tmpJson{Limit: 2, Hostname: "hiank"})
	assert.DeepEqual(t, yv, tmpYaml{A: "ws", Hope: "love", Hp: "hp", Slave: "slave"})

}
