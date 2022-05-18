package think

import "context"

var (
	Export_defaultOptions = func() Export_options {
		opts := defaultOptions()
		return Export_options{&opts}
	}
)

type Export_options struct {
	*options
}

func (eo *Export_options) NatsUrl() string {
	return eo.natsUrl
}

func (eo *Export_options) TODO() context.Context {
	return eo.todo
}

func (eo *Export_options) Mdb() map[DBTag]DB {
	return eo.mdb
}

func (eo *Export_options) Options() *options {
	return eo.options
}
