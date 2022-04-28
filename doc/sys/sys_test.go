package sys_test

import (
	"errors"
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

func TestFormat(t *testing.T) {
	f := sys.Export_formatFromPath("test.json")
	assert.Equal(t, f, sys.FormatJson)

	f = sys.Export_formatFromPath("test.JsOn")
	assert.Equal(t, f, sys.FormatJson, "not case sensitive")

	f = sys.Export_formatFromPath("test.yAML")
	assert.Equal(t, f, sys.FormatYaml)

	f = sys.Export_formatFromPath("test")
	assert.Equal(t, f, sys.FormatUnsupport)

	f = sys.Export_formatFromPath("test.xlsx")
	assert.Equal(t, f, sys.FormatUnsupport)
}

func TestBytes(t *testing.T) {
	b, err := sys.Export_formatoBytes(sys.FormatJson, func() ([]byte, error) { return []byte(tmpJsonValue), nil })
	assert.Equal(t, err, nil, err)
	var v tmpJson
	err = b.UnmarshalTo(&v)
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, v, tmpJson{Id: 9527, Name: "华安"})

	_, err = sys.Export_formatoBytes(sys.Format(11), func() ([]byte, error) { return []byte(tmpJsonValue), nil })
	assert.Equal(t, err, sys.ErrUnsupportFormat)

	_, err = sys.Export_formatoBytes(sys.FormatYaml, func() ([]byte, error) { return nil, errors.New("unimplement") })
	assert.Equal(t, err.Error(), "unimplement", err)

	b, err = sys.Export_formatoBytes(sys.FormatYaml, func() ([]byte, error) {
		return []byte(tmpJsonValue), nil
	})
	v.Id, v.Name = 11, "name"
	assert.Equal(t, err, nil, err)
	err = b.UnmarshalTo(&v)
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, v, tmpJson{Id: 11, Name: "name"})

	var yv tmpYaml
	err = b.UnmarshalTo(&yv)
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, yv, tmpYaml{})

	b, err = sys.Export_formatoBytes(sys.FormatYaml, func() ([]byte, error) { return []byte(tmpYamlValue), nil })
	assert.Equal(t, err, nil, err)
	b.UnmarshalTo(&yv)
	assert.DeepEqual(t, yv, tmpYaml{Key: "bbq", Value: 25})
}

func TestUnmarshal(t *testing.T) {
	v, err := sys.UnmarshalNew[*tmpJson]("testdata/config.json")
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, v, &tmpJson{Limit: 11, Hostname: "hiank"})

	err = sys.UnmarshalTo("testdata/dep/dep.json", v)
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, v, &tmpJson{Limit: 1, Hostname: "hiank"})
	// var v2 tmpYaml
	// v2, err := sys.Un

	v2, err := sys.UnmarshalNew[tmpYaml]("testdata/config.yaml")
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, v2, tmpYaml{A: "ws", Hope: "love", Hp: "hp", Slave: "slave"})

	err = sys.UnmarshalTo("testdata/dep/dep.YaMl", &v2)
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, v2, tmpYaml{A: "love-ws", Hope: "love", Hp: "hp", Slave: "slave"})

	_, err = sys.UnmarshalNew[[]tmpYaml]("testdata/config.yaml")
	assert.Assert(t, err != nil)

	_, err = sys.UnmarshalNew[string]("testdata/config.yaml")
	assert.Assert(t, err != nil)
}

func TestFat(t *testing.T) {
	fat := sys.NewFat()
	fat.LoadFiles("testdata/dep")
	var jv tmpJson
	var yv tmpYaml
	err := fat.UnmarshalTo(&jv, &yv)
	assert.Equal(t, err, nil, err)
	assert.DeepEqual(t, jv, tmpJson{Limit: 1})
	assert.DeepEqual(t, yv, tmpYaml{A: "love-ws"})

	fat.LoadFiles("testdata")
	fat.UnmarshalTo(&jv, &yv)
	//testdata/dep 不会再加载
	//testdata/dep2 最后加载
	assert.DeepEqual(t, jv, tmpJson{Limit: 2, Hostname: "hiank"})
	assert.DeepEqual(t, yv, tmpYaml{A: "ws", Hope: "love", Hp: "hp", Slave: "slave"})
}
