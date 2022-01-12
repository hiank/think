package config_test

import (
	"testing"

	"github.com/hiank/think/config"
	"gotest.tools/v3/assert"
)

type testConfig struct {
	Limit int    `json:"sys.Limit"`
	Key   string `yaml:"a"`
}

func TestUnmarshaler(t *testing.T) {
	// t.Run()
	u := config.NewParser()
	u.LoadFile("testdata", "testdata/config.json")

	var cfg testConfig
	u.ParseAndClear(&cfg)

	assert.Equal(t, cfg.Key, "love-ws")
	assert.Equal(t, cfg.Limit, 2)
}