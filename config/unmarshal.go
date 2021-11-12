package config

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

var (
	jsonSuffix = "json"
	yamlSuffix = "yaml"
)

func suffix(path string) (val string) {
	if idx := strings.LastIndexByte(path, '.'); idx != -1 {
		val = strings.ToLower(path[idx+1:])
	}
	return
}

type unmarshaler interface {
	unmarshal(...IConfig)
}

type jsonData struct {
	data []byte
}

func (jd *jsonData) unmarshal(vals ...IConfig) {
	for _, cfg := range vals {
		if err := json.Unmarshal(jd.data, cfg); err != nil {
			klog.Warningf("json.Unmarshal to %v: %v", cfg, err)
		}
	}
}

type yamlData struct {
	data []byte
}

func (yd *yamlData) unmarshal(vals ...IConfig) {
	for _, cfg := range vals {
		if err := yaml.Unmarshal(yd.data, cfg); err != nil {
			klog.Warningf("yaml.Unmarshal to %v: %v", cfg, err)
		}
	}
}

// type excelData struct {
// 	data []byte
// }

// func (ed *excelData) unmarshal(vals ...IConfig) {
// 	klog.Warning("not support excel now")
// }
