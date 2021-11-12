package config

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

type testConfig struct {
	Tik string `json:"sys.Tik"`
	Tok int    `yaml:"Tok"`
}

func TestJsonUnmarshaler(t *testing.T) {
	u := &jsonData{data: []byte(`{"sys.Tik": "nil"}`)}
	var cfg testConfig
	u.unmarshal(&cfg)

	assert.Equal(t, cfg.Tik, "nil")

	u = &jsonData{data: []byte(`{"sys.Tik": "overwrite"}`)}
	u.unmarshal(&cfg)
	assert.Equal(t, cfg.Tik, "overwrite", "overwrite previous settings")
}

func TestYamlUnmarshaler(t *testing.T) {
	u := &yamlData{data: []byte(`Tok: 2`)}
	var cfg testConfig
	u.unmarshal(&cfg)

	assert.Equal(t, cfg.Tok, 2)

	u = &yamlData{data: []byte(`Tok: 3`)}
	u.unmarshal(&cfg)
	assert.Equal(t, cfg.Tok, 3, "overwrite previous settings")
}

func TestMarch(t *testing.T) {
	paths := match("testdata")
	assert.Equal(t, len(paths), 5)

	root, _ := filepath.Abs("testdata")
	// strings.
	sp := string(filepath.Separator)
	root += sp
	// t.Log(root)
	names := []string{
		"config.json",
		"config.yaml",
		"dep" + sp + "dep.YaMl",
		"dep" + sp + "dep.json", //NOTE: 'j' > 'Y'
		"dep2" + sp + "dep.json",
	}
	for i, path := range paths {
		assert.Equal(t, path, root+names[i])
	}

	path := match("testdata/config.json")[0]
	assert.Equal(t, path, paths[0])

	paths = match("testdata/non.json")
	assert.Equal(t, len(paths), 0)

	_, err := ioutil.ReadDir("testdata")
	assert.Assert(t, err == nil)
}
