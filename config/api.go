package config

type IConfig interface{}

type IParser interface {
	LoadFile(paths ...string)
	LoadYamlBytes(values []byte)
	LoadJsonBytes(values []byte)

	ParseAndClear(configs ...IConfig)
}
