# Token and Token's Builder

`Builder` 用于构建存储`Token` 可以删除`Token` 单例对象，每个进程只能存在一次此对象(调用`Cancel` 后，整个进程中`token` 包失效，所以只在进程结束时调用`Cancel`)

`Token` 用于管理维护每个用户连接对应的资源

## API

- [Builder](#builder)
  - [Get](#get)
  - [Find](#find)
  - [Build](#build)
  - [Delete](#delete)
  - [Cancel](#cancel)

- [Token](#token)
  - [Derive](#derive)
  - [ToString](#derive)
  - [Cancel](#cancel)
  - [Context](#Context)

## Builder

### Get

**Param**: `string`

**Return**: `*Token`, `error`

用于通过`string` key 获取`Token`
- 如果`Token` 已存在，则直接返回
- 构建`Token`
  - 构建成果，存储并返回结果
  - 构建失败，返回错误

### Find

**Param**: `string`

**Return**: `*Token`, `bool`

用于通过`string` key 查找已构建的`Token`
- 找到，返回结果
- 未找到，返回`false`

### Build

**Parem**: `string`

**Return**: `*Token`, `error`

用于通过`string` key 构建`Token` 并存储
- 已存在`Token`，返回错误
- 不存在`Token`，构建`Token` 并存储，返回此`Token`


### Delete

**Parem**: `string`

用于通过`string` key 删除存储的`Token`

### Cancel

用于清除`Builder`，当进程结束时才可以调用此方法，调用后，单例`Builder` 将被置为`nil`，并且无法重新构建。

## Token

### Derive

**Return**: `*Token`

生成派生`Token`, 对于那些依赖与主`Token` 但不影响主`Token` 的资源，需要使用派生的`Token`。当此类资源需要释放时，只需要关闭派生的`Token` 即可，主`Token` 并不受影响。当主`Token` 被关闭时，派生`Token` 也会收到关闭信号

### ToString

**Return**: `string`

返回此`Token` 的`string` key

### Cancel

关闭，如果是主`Token` 则会从单例`Builder` 中删除之

### Context

设定为一个`context.Context` 对象，可以方便的使用`context.Context` 的相关方法
