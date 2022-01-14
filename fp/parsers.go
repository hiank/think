package fp

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

type form int8

const (
	formNon  form = 0
	formJson form = 1 << 0
	formYaml form = 1 << 1
	// formExcel form = 1 << 2
	// formXls   form = (1 << 3) | formExcel
	// formXlsx form = (1 << 3) | formExcel
)

func fileForm(path string) (f form) {
	if idx := strings.LastIndexByte(path, '.'); idx != -1 {
		switch strings.ToLower(path[idx+1:]) {
		case "json":
			f = formJson
		case "yaml":
			f = formYaml
		}
	}
	return
}

type parser interface {
	parse(...interface{})
}

type jsonData struct {
	data []byte
}

func (jd *jsonData) parse(vals ...interface{}) {
	for _, val := range vals {
		if err := json.Unmarshal(jd.data, val); err != nil {
			klog.Warningf("json.Unmarshal to %v: %v", val, err)
		}
	}
}

type yamlData struct {
	data []byte
}

func (yd *yamlData) parse(vals ...interface{}) {
	for _, val := range vals {
		if err := yaml.Unmarshal(yd.data, val); err != nil {
			klog.Warningf("yaml.Unmarshal to %v: %v", val, err)
		}
	}
}
