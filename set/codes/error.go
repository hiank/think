package codes

import (
	"errors"
)

// //ErrorNilValue 错误码，数据为nil
// var ErrorNilValue = errors.New("value is nil")

// //ErrorNotExisted 错误码，不存在
// var ErrorNotExisted = errors.New("not existed")

// //ErrorNotSupportType 错误码，不支持的类型
// var ErrorNotSupportType = errors.New("node support type")

// //ErrorNoMessageHandler 没有对应的handler用于处理此消息
// var ErrorNoMessageHandler = errors.New("no message handler for Handle the message")

// //ErrorExistedMessageHandler 已存在对应的handler
// var ErrorExistedMessageHandler = errors.New("existed message handler already")

// //ErrorAnyMessageIsEmpty
// var ErrorAnyMessageIsEmpty = errors.New("any message is empty")

const (
	ErrorUnknown               = 100
	ErrorNilValue              = 101
	ErrorNotExisted            = 102
	ErrorNotSupportType        = 103
	ErrorNoMessageHandler      = 104
	ErrorExistedMessageHandler = 105
	ErrorAnyMessageIsEmpty     = 106
	ErrorSrvClosed             = 107
)

var errorCache = map[int]error{
	ErrorUnknown:               errors.New("unknown error"),
	ErrorNilValue:              errors.New("value is nil"),
	ErrorNotExisted:            errors.New("not existed"),
	ErrorNotSupportType:        errors.New("node support type"),
	ErrorNoMessageHandler:      errors.New("no message handler for Handle the message"),
	ErrorExistedMessageHandler: errors.New("existed message handler already"),
	ErrorAnyMessageIsEmpty:     errors.New("any message is empty"),
	ErrorSrvClosed:             errors.New("srv closed"),
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
