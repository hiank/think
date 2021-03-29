package codes

const (
	PanicNilHandler     = 101
	PanicExistedHandler = 102
	PanicNonLimit       = 103
)

var panicText = map[int]string{
	PanicNilHandler:     "handler can not be nil",
	PanicExistedHandler: "handler already existed",
	PanicNonLimit:       "limitMux max can not be 0",
}

// PanicText returns a text for the panic code. It returns the empty
// string if the code is unknown.
func PanicText(code int) string {
	return panicText[code]
}
