package codes

import "errors"

//ErrorNilValue 错误码，数据为nil
var ErrorNilValue = errors.New("value is nil")

//ErrorNotExisted 错误码，不存在
var ErrorNotExisted = errors.New("not existed")

//ErrorNotSupportType 错误码，不支持的类型
var ErrorNotSupportType = errors.New("node support type")
