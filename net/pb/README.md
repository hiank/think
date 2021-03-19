# pb package

基础消息及其处理方法合集。服务会将收到的消息自动包装成pb.Message，加入一些固定识别信息。
提供一个简单的方式，可以方便的注册消息处理方法，注册后的方法将截断服务收到的消息。

## `pb.Message`

封装收到的消息，添加必要识别信息，用于后续识别处理。用于直接传递的消息务必包含特定头信息，可识别为不同的消息类型。

### 消息类型

- `TypeUndefined`: 未定义类型，将无法处理
- `TypeGET`: `GET`类型，以`G`开头，发送请求，获得返回，一次
- `TypePOST`: `POST`类型，以`P`开头，发送消息，无需返回
- `TypeSTREAM`: `STREAM`类型，以`S`开头，建立管道，消息自由收发
- `TypeMQ`: `MQ`类型，以`M`开头，使用消息中间件传递

### `GetServeType` 获取[消息类型](#消息类型)

```golang
GetServeType(anyMsg *anypb.Any) (t int, err error)
```

## `pb.LiteHandler`

主要提供以下功能：

- `Register` 注册某个协议的处理方法，服务方收到消息后将调用对应的Handler
- `Handler` 实现pool.Handler接口。Handler将优先使用注册的Handler来处理消息。如无注册Handler，则尝试调用DefaultHandler
- `DefaultHandler` 通用的Handler，在找不到注册Handler情况下，将尝试调用此方法

_NOTE:_ `非线程安全，务必在初始化阶段注册所有需要的Handler及DefaultHandler`
