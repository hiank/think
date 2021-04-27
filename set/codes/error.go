package codes

import (
	"errors"
)

const (
	ErrorUnknown                  = 100
	ErrorNilValue                 = 101
	ErrorWorkerNotExisted         = 102
	ErrorNotSupportType           = 103
	ErrorNoMessageHandler         = 104
	ErrorExistedMessageHandler    = 105
	ErrorAnyMessageIsEmpty        = 106
	ErrorNeedOneofConnRecvHandler = 107
	ErrorNonSupportLiteServe      = 108
	ErrorNonHelper                = 109
)

var errorCache = map[int]error{
	ErrorUnknown:                  errors.New("unknown error"),
	ErrorNilValue:                 errors.New("value is nil"),
	ErrorWorkerNotExisted:         errors.New("worker not existed, cannot operate the message"),
	ErrorNotSupportType:           errors.New("node support type"),
	ErrorNoMessageHandler:         errors.New("no message handler for Handle the message"),
	ErrorExistedMessageHandler:    errors.New("existed message handler already"),
	ErrorAnyMessageIsEmpty:        errors.New("any message is empty"),
	ErrorNeedOneofConnRecvHandler: errors.New("must exist connHandler or recvHandler"),
	ErrorNonSupportLiteServe:      errors.New("LiteSender depend for a liteServer, but not the server"),
	ErrorNonHelper:                errors.New("ListenAndServe need useful Helper"),
}

// Error returns an error value for the error code. It returns "unknow error"
// if the code is unknown.
func Error(code int) error {
	err, ok := errorCache[code]
	if !ok {
		err = errorCache[ErrorUnknown]
	}
	return err
}
