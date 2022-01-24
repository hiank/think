package net

// dialOptions configure a Dial call. dialOptions are set by the DialOption
// values passed to Dial.
type handleOptions struct {
	defaultHandler CarrierHandler
	converter      CarrierConverter
}

// HandleOption configures how we set up the connection.
type HandleOption interface {
	apply(*handleOptions)
}

// funcDialOption wraps a function that modifies dialOptions into an
// implementation of the DialOption interface.
type funcHandleOption func(*handleOptions)

func (fho funcHandleOption) apply(do *handleOptions) {
	fho(do)
}

func WithDefaultHandler(handler CarrierHandler) HandleOption {
	return funcHandleOption(func(ho *handleOptions) {
		ho.defaultHandler = handler
	})
}

func WithConverter(converter CarrierConverter) HandleOption {
	return funcHandleOption(func(ho *handleOptions) {
		ho.converter = converter
	})
}
