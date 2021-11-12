package config

type IConfig interface{}

type IUnmarshaler interface {
	HandleFile(paths ...string)
	HandleYamlBytes(values []byte)
	HandleJsonBytes(values []byte)

	UnmarshalAndClean(configs ...IConfig)
}
