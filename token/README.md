# Token and Token's Builder

`Builder` 用于构建存储`Token` 可以删除`Token` 单例对象，每个进程只能存在一次此对象(调用`Cancel` 后，整个进程中`token` 包失效，所以只在进程结束时调用`Cancel`)

`Token` 用于管理维护每个用户连接对应的资源

## API

- [Builder](#builder)
  - [Find](#find)
  - [Get](#get)
  - [removeReq](#removeReq)

- [Token](#token)
  - [Derive](#derive)
  - [ToString](#derive)
  - [Cancel](#cancel)
  - [Context](#Context)

## Builder

- `Builder`自构建始，整个程序活动期间都将存在，不会被消除。当前版本只支持两个方法，Get 和 removeReq
- 使用chan 带缓存方式避免调用goroutine 阻塞，缓存大小可配置

### Find
- **Param**: `string`
- **Return**: `*Token`, `bool`

用于判断及查找主`Token`，如果主`Token`已经找不到了，可能要丢弃处理。否则可能会派生一个`Token`与相关资源绑定

### Get

- **Param**: `string`
- **Return**: `*Token`

用于通过`string` key 获取`Token`，使用通信的方式获取值，避免数据竞争[为了与`removeReq`表现一致，没有使用读写锁，而使用了chan 来同步]
当前的方法必然能返回一个`*Token`

### removeReq

- **Return**: `chan<- *Token`

删除管道写入的*Token. 避免大量token失效时，请求被等待而挂起

## Token

分为主`Token`和派生`Token`，主`Token`与用户切实相关，所以只会在用于连接建立后才构建此`Token`，其余场景可直接引用或使用派生的`Token`
1. 主`Token`:
    - 受`Builder`维护管理，失效后须通知`Builder`以删除引用[`Cancel`方法将执行此操作]
    - 有超时处理，超时后将自调用`Cancel`方法进行清理
    - 可派生子`Token`，特性参加下述说明
    - 提供清理函数`Cancel`，将关闭关联的Context，以广播通知所有监听此`Token`状态的相关方法；所有派生的`Token`将一并收到此关闭消息
2. 派生`Token`:
    - 由主`Token`或其它派生`Token`调用`Derive`方法构建得来
    - 生命周期上限为派生出此派生`Token`的`Token`
    - 不执行超时监听
    - 对此`Token`的关闭操作不影响父`Token` 

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
