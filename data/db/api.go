package db

//IClient database client
type IClient interface {
	HGet(hashKey, fieldKey string) (IParser, error)

	// HSet accepts values in following formats:
	//   - HSet("myhash", "key1", "value1", "key2", "value2")
	//   - HSet("myhash", []string{"key1", "value1", "key2", "value2"})
	//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
	//
	// Note that it requires Redis v4 for multiple field/value pairs support.
	//refer to redis Client.HSet
	HSet(hashKey string, values ...interface{}) error

	Close() error
}

//IParser value parser
//refer to redis
type IParser interface {
	Scan(interface{}) error
	Bool() (bool, error)
	Int() (int, error)
	Int64() (int64, error)
	Uint64() (uint64, error)
	Result() (string, error)
	Float32() (float32, error)
	Float64() (float64, error)
}
