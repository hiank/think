package codes

import "errors"

//ErrorNilValue 错误码，数据为nil
var ErrorNilValue = errors.New("value is nil")

//ErrorNotExisted 错误码，不存在
var ErrorNotExisted = errors.New("not existed")

//ErrorNotSupportType 错误码，不支持的类型
var ErrorNotSupportType = errors.New("node support type")

//ErrorNoMessageHandler 没有对应的handler用于处理此消息
var ErrorNoMessageHandler = errors.New("no message handler for Handle the message")

//ErrorExistedMessageHandler 已存在对应的handler
var ErrorExistedMessageHandler = errors.New("existed message handler already")

//ErrorAnyMessageIsEmpty
var ErrorAnyMessageIsEmpty = errors.New("any message is empty")
