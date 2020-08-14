# Message

消息数据，protobuf格式，用于think中传导

## Type

消息类型，每个消息名须以`G || P || S`开头，对应不同的处理，这个主要是k8s中需要区分不同的消息

- `TypeUndefined`: 未定义类型，将无法处理
- `TypeGET`: `GET`类型，消息名以`G`开头，发送请求，获得返回，一次
- `TypePOST`: `POST`类型，消息名以`P`开头，发送消息，无需返回
- `TypeSTREAM`: `STREAM`类型，消息名以`S`开头，建立管道，消息自由收发


## API

- [AnyDecode([]byte) (*any.Any, error)](#AnyDecode)
- [AnyEncode(*any.Any) ([]byte, error)](#AnyEncode)
- [AnyMessageNameTrimed(*any.Any) (string, error)](#AnyMessageNameTrimed)
- [GetServerType(*any.Any) (int, error)](#GetServerType)


## `AnyDecode`

将[]byte 消息编码为Any 消息

## `AnyEncode`

将Any 消息编码为[]byte 消息

## `AnyMessageNameTrimed`

获取处理过的any.Any 消息名，去掉可能包含的报名

## `GetServerType`

获取服务类型，用于使用不同的方式处理消息